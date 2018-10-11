package formatter

import (
	"log"
	"os"
	"reflect"
	"testing"
)

var fakeLogger = log.New(&mockWriter{}, "tests: ", log.LstdFlags)

func TestNewFormatter(t *testing.T) {
	type args struct {
		t Type
	}
	tests := []struct {
		name string
		args args
		want Formatter
	}{
		{"1", args{t: StdoutType}, &stdout{l: fakeLogger}},
		{"2", args{t: FileType}, &file{l: fakeLogger}},
		{"3", args{t: JSONType}, &json{l: fakeLogger}},
		{"4", args{t: YamlType}, &yaml{l: fakeLogger}},
		{"5", args{t: 33}, &stdout{l: fakeLogger}},
		{"6", args{t: SkopeoType}, &skopeo{l: fakeLogger}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFormatter(tt.args.t, "", fakeLogger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockWriter struct {
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return 0, nil
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
		{"1", args{"t.txt", []byte("test content"), fakeLogger}, false},
		{"2", args{"", []byte("test content"), fakeLogger}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.name, tt.args.content, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	os.RemoveAll("t.txt")
}
