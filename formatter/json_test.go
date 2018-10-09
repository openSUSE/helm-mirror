package formatter

import (
	"bytes"
	"os"
	"testing"
)

func Test_json_Output(t *testing.T) {
	var buff bytes.Buffer
	buff.WriteString("test")
	var buffError bytes.Buffer
	buffError.WriteString(`?est: ]`)
	type args struct {
		b bytes.Buffer
	}
	tests := []struct {
		name    string
		f       *json
		args    args
		wantErr bool
	}{
		{"1", &json{fileName: "test.json", l: fakeLogger}, args{buff}, false},
		{"2", &json{fileName: "", l: fakeLogger}, args{buff}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Output(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("json.Output() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("test.json")
}
