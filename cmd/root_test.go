package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/openSUSE/helm-mirror/fixtures"
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
	svr := fixtures.StartHTTPServer()
	defer svr.Shutdown(nil)
	fixtures.WaitForServer("http://127.0.0.1:1793/alive")
	type args struct {
		cmd          *cobra.Command
		args         []string
		newRootURL   string
		ignoreErrors bool
		allVersions  bool
		chartName    string
		chartVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{&cobra.Command{}, []string{"http://test", path.Join("/mr", "mzxyptlk")}, "", false, true, "", ""}, true},
		{"2", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "", true, true, "", ""}, false},
		{"3", args{&cobra.Command{}, []string{"%", dir}, "", false, true, "", ""}, true},
		{"4", args{&cobra.Command{}, []string{"http://test", dir}, "%", false, true, "", ""}, true},
		{"5", args{&cobra.Command{}, []string{"http://test", dir}, "ftp://test", false, true, "", ""}, true},
		{"6", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "https://test/com/charts", true, true, "", ""}, false},
		{"7", args{&cobra.Command{}, []string{"http://127.0.0.1:1111", dir}, "https://test/com/charts", false, true, "", ""}, true},
		{"8", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "https://test/com/charts", true, true, "", ""}, false},
		{"9", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "", false, true, "", ""}, true},
		{"10", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "", true, false, "", ""}, false},
		{"11", args{&cobra.Command{}, []string{"http://127.0.0.1:1793", dir}, "", true, false, "", "1.0.0"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("%s clean\n", filepath.Clean(dir))
			fmt.Printf("%v\n", tt.args.args[1])
			newRootURL = tt.args.newRootURL
			IgnoreErrors = tt.args.ignoreErrors
			AllVersions = tt.args.allVersions
			chartName = tt.args.chartName
			chartVersion = tt.args.chartVersion
			if err := runRoot(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
