package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runVersion(t *testing.T) {
	type args struct {
		in0 *cobra.Command
		in1 []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{&cobra.Command{}, []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runVersion(tt.args.in0, tt.args.in1)
		})
	}
}
