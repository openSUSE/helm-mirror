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
	gService := &GetService{config: config, logger: fakeLogger, newRootURL: "https://newchartserver.com", allVersions: false}
	type args struct {
		helmRepo     string
		workspace    string
		verbose      bool
		ignoreErrors bool
		logger       *log.Logger
		newRootURL   string
		allVersions  bool
		chartName    string
		chartVersion string
	}
	tests := []struct {
		name string
		args args
		want GetServiceInterface
	}{
		{"1", args{"http://helmrepo", dir, false, false, fakeLogger, "https://newchartserver.com", false, "", ""}, gService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetService(config, tt.args.verbose, tt.args.allVersions, tt.args.ignoreErrors, tt.args.logger, tt.args.newRootURL, tt.args.chartName, tt.args.chartVersion); !reflect.DeepEqual(got, tt.want) {
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
	defer svr.Shutdown(nil)
	fixtures.WaitForServer("http://127.0.0.1:1793/alive")
	type fields struct {
		repoURL      string
		workDir      string
		ignoreErrors bool
		verbose      bool
		allVersions  bool
		chartName    string
		chartVersion string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantTgz int
	}{
		{"1", fields{"", "", false, false, true, "", ""}, true, 0},
		{"2", fields{"http://127.0.0.1", "", false, false, true, "", ""}, true, 0},
		{"3", fields{"http://127.0.0.1:1793", "", false, false, true, "", ""}, true, 0},
		{"4", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), false, false, true, "", ""}, true, 0},
		{"5", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, false, true, "", ""}, false, 4},
		{"6", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, false, false, "", ""}, false, 3},
		{"7", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, "", ""}, false, 3},
		{"8", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, "chart2", ""}, false, 1},
		{"9", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, "chart", ""}, false, 0},
		{"10", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, `^(?:(?:aa)|.$`, ""}, true, 0},
		{"11", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, "chart2", "7.0.0"}, false, 0},
		{"12", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, false, "chart2", "0.0.0-rc1"}, false, 1},
		{"13", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, true, "chart2", ""}, false, 2},
		{"14", fields{"http://127.0.0.1:1793", path.Join(dir, "get"), true, true, true, "chart2", "0.0.0-rc1"}, false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GetService{
				config:       repo.Entry{Name: tt.fields.workDir, URL: tt.fields.repoURL},
				logger:       fakeLogger,
				ignoreErrors: tt.fields.ignoreErrors,
				verbose:      tt.fields.verbose,
				allVersions:  tt.fields.allVersions,
				chartName:    tt.fields.chartName,
				chartVersion: tt.fields.chartVersion,
			}
			if err := g.Get(); (err != nil) != tt.wantErr {
				t.Errorf("GetService.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				files, err := ioutil.ReadDir(path.Join(dir, "get"))
				if err != nil {
					log.Fatal(err)
				}
				count := 0
				for _, f := range files {
					if strings.Contains(f.Name(), ".tgz") {
						count++
					}
					os.RemoveAll(path.Join(dir, "get", f.Name()))
				}
				if count != tt.wantTgz {
					t.Errorf("GetService.Get() got count of = %v TGZ files, want count of %v", count, tt.wantTgz)
				}
			}
		})
	}
	os.RemoveAll("downloaded-index.yaml")
}

func Test_writeFile(t *testing.T) {
	type args struct {
		name         string
		content      []byte
		log          *log.Logger
		ignoreErrors bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{"tmp.txt", []byte("test"), fakeLogger, false}, false},
		{"2", args{"", []byte("test"), fakeLogger, false}, true},
		{"2", args{"", []byte("test"), fakeLogger, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.name, tt.args.content, tt.args.log, tt.args.ignoreErrors); (err != nil) != tt.wantErr {
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
		folder       string
		URL          string
		newRootURL   string
		log          *log.Logger
		ignoreErrors bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{path.Join(dir, "processfolder"), "http://127.0.0.1:1793", "http://newchart.server.com", fakeLogger, false}, false},
		{"2", args{path.Join(dir, "processerrorfolder"), "http://127.0.0.1:1793", "http://newchart.server.com", fakeLogger, false}, true},
		{"3", args{path.Join(dir, "processfolder"), "http://127.0.0.1:1793", "", fakeLogger, false}, false},
	}
	for _, tt := range tests {
		ioutil.WriteFile(path.Join(dir, "processfolder", "downloaded-index.yaml"), []byte(fixtures.IndexYaml), 0666)
		t.Run(tt.name, func(t *testing.T) {
			if err := prepareIndexFile(tt.args.folder, tt.args.URL, tt.args.newRootURL, tt.args.log, tt.args.ignoreErrors); (err != nil) != tt.wantErr {
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
					t.Errorf("prepareIndexFile() replacedCount = %v, want replacedCount %v", count, fixtures.Expectedcharts)
				}
				_, err = os.Stat(path.Join(dir, "processfolder", "downloaded-index.yaml"))
				if err == nil {
					t.Errorf("prepareIndexFile() dowloaded-index.yaml not deleted")
				}
			}
		})
	}
}
