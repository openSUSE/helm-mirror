package cmd

import (
	"log"
	"os"
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
	type args struct {
		output   string
		fileName string
		l        *log.Logger
	}
	tests := []struct {
		name string
		args args
		want formatter.Formatter
	}{
		{"1", args{"stdout", "test", fakeLog}, formatter.NewFormatter(formatter.StdoutType, "test", fakeLog)},
		{"2", args{"file", "/test", fakeLog}, formatter.NewFormatter(formatter.FileType, "/test", fakeLog)},
		{"3", args{"yaml", "/test", fakeLog}, formatter.NewFormatter(formatter.YamlType, "/test", fakeLog)},
		{"4", args{"json", "/test", fakeLog}, formatter.NewFormatter(formatter.JSONType, "/test", fakeLog)},
		{"5", args{"notexists", "test", fakeLog}, formatter.NewFormatter(formatter.StdoutType, "test", fakeLog)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveFormatter(tt.args.output, tt.args.fileName, tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runInspectImages(t *testing.T) {
	var cmd = &cobra.Command{}
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", outputDesc)
	cmd.PersistentFlags().StringVar(&imagesFile, "file-name", "images.out", fileNameDesc)
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
