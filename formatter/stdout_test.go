package formatter

import (
	"bytes"
	"testing"
)

func Test_stdout_Output(t *testing.T) {
	var buff bytes.Buffer
	buff.WriteString("test")
	type args struct {
		b bytes.Buffer
	}
	tests := []struct {
		name    string
		s       *stdout
		args    args
		wantErr bool
	}{
		{"1", &stdout{l: fakeLogger}, args{buff}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Output(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("stdout.Output() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
