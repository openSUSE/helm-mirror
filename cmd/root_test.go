package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func Test_validateRootArgs(t *testing.T) {
	c := &cobra.Command{}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{c, []string{}}, true},
		{"2", args{c, []string{"url"}}, true},
		{"3", args{c, []string{"url", "target"}}, true},
		{"4.1", args{c, []string{"http://url", "target"}}, true},
		{"4.2", args{c, []string{"https://url", "target"}}, true},
		{"4.3", args{c, []string{"ftp://url", "target"}}, true},
		{"4.4", args{c, []string{"ftps://url", "target"}}, true},
		{"5.1", args{c, []string{"http://url", "/target"}}, false},
		{"5.2", args{c, []string{"https://url", "/target"}}, false},
		{"5.3", args{c, []string{"ftp://url", "/target"}}, true},
		{"5.4", args{c, []string{"ftps://url", "/target"}}, true},
		{"6", args{c, []string{"ftps://url", "/target", "extra"}}, true},
		{"7", args{c, []string{"help"}}, false},
		{"8", args{c, []string{"%", "/target", "extra"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRootArgs(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validateRootArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runRoot(t *testing.T) {
	dir, err := ioutil.TempDir("", "helmmirror")
	if err != nil {
		t.Errorf("creating temp dir: %s", err)
	}
	defer os.RemoveAll(dir)
	svr := startHTTPServer()
	type args struct {
		cmd        *cobra.Command
		args       []string
		newRootURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{&cobra.Command{}, []string{"http://test", path.Join("/mr", "mzxyptlk")}, ""}, true},
		{"2", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, ""}, false},
		{"3", args{&cobra.Command{}, []string{"%", dir}, ""}, true},
		{"4", args{&cobra.Command{}, []string{"http://test", dir}, "%"}, true},
		{"5", args{&cobra.Command{}, []string{"http://test", dir}, "ftp://test"}, true},
		{"6", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "https://test/com/charts"}, false},
		{"7", args{&cobra.Command{}, []string{"http://127.0.0.1:1111", dir}, "https://test/com/charts"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("%s clean\n", filepath.Clean(dir))
			fmt.Printf("%v\n", tt.args.args[1])
			newRootURL = tt.args.newRootURL
			if err := runRoot(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if err := svr.Shutdown(nil); err != nil {
		t.Errorf("error stoping down web server")
	}
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
	w.Write(indexYaml)
}

func chartTgz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "binary/octet-stream")
	w.Write(chartTGZ)
}

var indexYaml = []byte(`apiVersion: v1
entries:
  chart1:
  - apiVersion: v2
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing
    digest: 8cc99f9cb669171776f7c6ec66069907579be91179f9201725fc6fc6f9ef1f29
    name: chart1
    urls:
    - http://127.0.0.1:1793/chart1-2.11.0.tgz
    version: 2.11.0
  chart2:
  - apiVersion: v1
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing too
    digest: 0c76ee9b4b78cb60fcce8c00ec0f5048cbe626fcaabe48f2f8e84b029e894f49
    name: chart2
    urls:
    - http://127.0.0.1:1793/chart2-1.0.1.tgz
    version: 1.0.1
  - apiVersion: v1
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing too
    digest: e9a545006570b7fc5e4458f6eae178c2aa8f8e9e57eafac59869c856b86e862f
    name: chart2
    urls:
    - http://127.0.0.1:1793/chart2-0.0.0-pre.tgz
    version: 0.0.0-pre
`)

var chartTGZ = []byte{31, 139, 8, 0, 224, 223, 181, 91, 0, 3, 237, 193, 1, 13, 0, 0, 0, 194,
	160, 247, 79, 109, 14, 55, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 55, 3, 154, 222, 29, 39, 0, 40, 0, 0}
