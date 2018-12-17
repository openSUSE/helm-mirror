package service

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/openSUSE/helm-mirror/fixtures"

	"k8s.io/helm/pkg/repo"
)

var fakeLogger = log.New(&mockLog{}, "test:", log.LstdFlags)

type mockLog struct{}

func (m *mockLog) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestNewGetService(t *testing.T) {
	dir, err := ioutil.TempDir("", "helmmirrortests")
	if err != nil {
		t.Errorf("Creating tmp directory: %s", err)
	}
	defer os.RemoveAll(dir)
	config := repo.Entry{Name: dir, URL: "http://helmrepo"}
	gService := &GetService{config: config, logger: fakeLogger, newRootURL: "https://newchartserver.com"}
	type args struct {
		helmRepo     string
		workspace    string
		verbose      bool
		ignoreErrors bool
		logger       *log.Logger
		newRootURL   string
	}
	tests := []struct {
		name string
		args args
		want GetServiceInterface
	}{
		{"1", args{"http://helmrepo", dir, false, false, fakeLogger, "https://newchartserver.com"}, gService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetService(config, tt.args.verbose, tt.args.ignoreErrors, tt.args.logger, tt.args.newRootURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService_Get(t *testing.T) {
	dir, err := prepareTmp()
	if err != nil {
		t.Errorf("loading testdata: %s", err)
	}
	defer os.RemoveAll(dir)
	svr := fixtures.StartHTTPServer()
	fixtures.WaitForServer(svr.Addr)
	type fields struct {
		repoURL      string
		workDir      string
		ignoreErrors bool
		verbose      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"1", fields{"", "", false, false}, true},
		{"2", fields{"http://127.0.0.1", "", false, false}, true},
		{"3", fields{"http://127.0.0.1:1793", "", false, false}, true},
		{"4", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), false, false}, true},
		{"5", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetService{
				config:       repo.Entry{Name: tt.fields.workDir, URL: tt.fields.repoURL},
				logger:       fakeLogger,
				ignoreErrors: tt.fields.ignoreErrors,
				verbose:      tt.fields.verbose,
			}
			if err := g.Get(); (err != nil) != tt.wantErr {
				t.Errorf("GetService.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("downloaded-index.yaml")
	if err := svr.Shutdown(nil); err != nil {
		t.Errorf("error stoping down web server")
	}
}

func Test_writeFile(t *testing.T) {
	type args struct {
		name    string
		content []byte
		log     *log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"tmp.txt", []byte("test"), fakeLogger}, false},
		{"2", args{"", []byte("test"), fakeLogger}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.name, tt.args.content, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("tmp.txt")
}

func Test_prepareIndexFile(t *testing.T) {
	dir, err := prepareTmp()
	if err != nil {
		t.Errorf("loading testdata: %s", err)
	}
	defer os.RemoveAll(dir)
	type args struct {
		folder     string
		URL        string
		newRootURL string
		log        *log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{path.Join(dir, "processfolder"), "http://127.0.0.1:1793", "http://newchart.server.com", fakeLogger}, false},
		{"2", args{path.Join(dir, "processerrorfolder"), "http://127.0.0.1:1793", "http://newchart.server.com", fakeLogger}, true},
		{"3", args{path.Join(dir, "processfolder"), "http://127.0.0.1:1793", "", fakeLogger}, false},
	}
	for _, tt := range tests {
		ioutil.WriteFile(path.Join(dir, "processfolder", "downloaded-index.yaml"), []byte(fixtures.IndexYaml), 0666)
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareIndexFile(tt.args.folder, tt.args.URL, tt.args.newRootURL, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("prepareIndexFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "1" {
				contentBytes, err := ioutil.ReadFile(path.Join(dir, "processfolder", "index.yaml"))
				if err != nil {
					t.Log("Error reading index.yaml")
				}
				content := string(contentBytes)
				count := strings.Count(content, tt.args.newRootURL)
				if count != fixtures.Expectedcharts {
					t.Errorf("prepareIndexFile() replacedCount = %v, want replacedCount %v", count, 3)
				}
				_, err = os.Stat(path.Join(dir, "processfolder", "downloaded-index.yaml"))
				if err == nil {
					t.Errorf("prepareIndexFile() dowloaded-index.yaml not deleted")
				}
			}
		})
	}
}
