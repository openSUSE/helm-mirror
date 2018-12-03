package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/openSUSE/helm-mirror/formatter"
	"github.com/spf13/cobra"
)

func Test_validateInspectImagesArgs(t *testing.T) {
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
		{"2", args{c, []string{"target"}}, true},
		{"3", args{c, []string{"target", "destination"}}, true},
		{"4", args{c, []string{"/target", "destination"}}, false},
		{"5", args{c, []string{"/target/tar.tgz", "destination"}}, false},
		{"6", args{c, []string{"/target", "/destination", "extra"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInspectImagesArgs(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("validateInspectImagesArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resolveFormatter(t *testing.T) {
	abs, err := filepath.Abs("")
	if err != nil {
		t.Errorf("resolvefilePath() = %s", abs)
	}
	resultPath := path.Join(abs, "images.out")
	type args struct {
		output string
		l      *log.Logger
	}
	tests := []struct {
		name string
		args args
		want formatter.Formatter
	}{
		{"1", args{"stdout", fakeLog}, formatter.NewFormatter(formatter.StdoutType, "images.out", fakeLog)},
		{"2.1", args{"file", fakeLog}, formatter.NewFormatter(formatter.FileType, resultPath, fakeLog)},
		{"2.2", args{"file=/test.txt", fakeLog}, formatter.NewFormatter(formatter.FileType, "/test.txt", fakeLog)},
		{"3.1", args{"yaml", fakeLog}, formatter.NewFormatter(formatter.YamlType, resultPath, fakeLog)},
		{"3.2", args{"yaml=/test.yaml", fakeLog}, formatter.NewFormatter(formatter.YamlType, "/test.yaml", fakeLog)},
		{"4.1", args{"json", fakeLog}, formatter.NewFormatter(formatter.JSONType, resultPath, fakeLog)},
		{"4.2", args{"json=/test.json", fakeLog}, formatter.NewFormatter(formatter.JSONType, "/test.json", fakeLog)},
		{"5", args{"notexists", fakeLog}, formatter.NewFormatter(formatter.StdoutType, resultPath, fakeLog)},
		{"6", args{"skopeo=/skopeo.yaml", fakeLog}, formatter.NewFormatter(formatter.SkopeoType, "/skopeo.yaml", fakeLog)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := resolveFormatter(tt.args.output, tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runInspectImages(t *testing.T) {
	var cmd = &cobra.Command{}
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", outputDesc)
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{cmd, []string{"/tmp/target"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runInspectImages(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runInspectImages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("/tmp/target")
}

var fakeLog = log.New(&mockLog{}, "tests:", log.LstdFlags)

type mockLog struct{}

func (m *mockLog) Write(p []byte) (n int, err error) {
	return 0, nil
}
