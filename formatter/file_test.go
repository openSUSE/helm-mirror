package formatter

import (
	"bytes"
	"os"
	"testing"
)

func Test_file_Output(t *testing.T) {
	var buff bytes.Buffer
	buff.WriteString("test")
	type args struct {
		b bytes.Buffer
	}
	tests := []struct {
		name    string
		f       *file
		args    args
		wantErr bool
	}{
		{"1", &file{fileName: "test.txt", l: fakeLogger}, args{buff}, false},
		{"2", &file{fileName: "", l: fakeLogger}, args{buff}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Output(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("file.Output() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("test.txt")
}
