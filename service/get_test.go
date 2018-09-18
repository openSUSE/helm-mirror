package service

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	"k8s.io/helm/pkg/repo"
)

var fakeLogger = log.New(&mockLog{}, "test:", log.LstdFlags)

type mockLog struct{}

func (m *mockLog) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestNewGetService(t *testing.T) {
	config := repo.Entry{Name: "/tmp/helmmirrortest", URL: "http://helmrepo"}
	gService := &GetService{config: config, logger: fakeLogger}
	type args struct {
		helmRepo     string
		workspace    string
		verbose      bool
		ignoreErrors bool
		logger       *log.Logger
	}
	tests := []struct {
		name string
		args args
		want GetServiceInterface
	}{
		{"1", args{"http://helmrepo", "/tmp/helmmirrortest", false, false, fakeLogger}, gService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetService(config, tt.args.verbose, tt.args.ignoreErrors, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetService() = %v, want %v", got, tt.want)
			}
		})
	}
	os.RemoveAll("/tmp/helmmirrortest")
}

func TestGetService_Get(t *testing.T) {
	prepareTmp()
	svr := startHTTPServer()
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
		{"4", fields{"http://127.0.0.1:1793", tmp + "/get", false, false}, false},
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
	tearDownTmp()
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

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":1793"}
	http.HandleFunc("/index.yaml", indexFile)
	http.HandleFunc("/chart1-2.11.0.tgz", chartTgz)
	http.HandleFunc("/chart2-1.0.1.tgz", chartTgz)
	http.HandleFunc("/chart2-0.0.0-pre.tgz", chartTgz)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	return srv
}

func indexFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "binary/octet-stream")
	w.Write([]byte(indexYaml))
}

func chartTgz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "binary/octet-stream")
	w.Write([]byte(chartTGZ))
}
