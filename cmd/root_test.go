package cmd

import (
	"os"
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
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"1", args{&cobra.Command{}, []string{"http://test", "/mr/mzxyptlk"}}, true},
		{"2", args{&cobra.Command{}, []string{"http://test", "/tmp/helm-mirror/test"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runRoot(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("/mr/mzxyptlk")
	os.RemoveAll("/tmp/helm-mirror/test")
}
