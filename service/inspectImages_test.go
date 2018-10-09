package service

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/openSUSE/helm-mirror/formatter"
)

var buff bytes.Buffer
var tmp = "/tmp/mirror"
var fakeFormatter = &mockFormatter{}

func TestNewImagesService(t *testing.T) {
	type args struct {
		target    string
		formatter formatter.Formatter
	}
	tests := []struct {
		name string
		args args
		want ImagesServiceInterface
	}{
		{"1", args{"/folder", fakeFormatter}, &ImagesService{target: "/folder", formatter: fakeFormatter, logger: fakeLogger, buffer: buff, verbose: false, ignoreErrors: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImagesService(tt.args.target, false, false, fakeFormatter, fakeLogger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImagesService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImagesService_Images(t *testing.T) {
	prepareTmp()
	defer tearDownTmp()
	type fields struct {
		target       string
		buff         bytes.Buffer
		ignoreErrors bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"1", fields{tmp + "/processfolder", buff, false}, false},
		{"2", fields{tmp + "/testdata/chart6", buff, false}, true},
		{"3", fields{tmp + "/testdata/chart1.tgz", buff, false}, false},
		{"4", fields{tmp + "/mr/mzxyptlk", buff, false}, true},
		{"5", fields{tmp + "/processfoldererror", buff, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ImagesService{
				target:       tt.fields.target,
				formatter:    fakeFormatter,
				logger:       fakeLogger,
				buffer:       tt.fields.buff,
				ignoreErrors: tt.fields.ignoreErrors,
				verbose:      false,
			}
			if err := i.Images(); (err != nil) != tt.wantErr {
				t.Errorf("ImagesService.Images() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImagesService_processDirectory(t *testing.T) {
	prepareTmp()
	type fields struct {
		target    string
		imageFile string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"1", fields{tmp + "/processfolder", "images"}, false},
		{"2", fields{tmp + "/testdata", "images"}, true},
		{"3", fields{tmp + "/testdata/chart1.tgz", "images"}, true},
		{"4", fields{tmp + "/processfolder/chart6", "images"}, true},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		t.Run(tt.name, func(t *testing.T) {
			i := &ImagesService{
				target:    tt.fields.target,
				formatter: fakeFormatter,
				logger:    fakeLogger,
				buffer:    buf,
			}
			if err := i.processDirectory(tt.fields.target); (err != nil) != tt.wantErr {
				t.Errorf("ImagesService.processDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	tearDownTmp()
}

func TestImagesService_processTarget(t *testing.T) {
	prepareTmp()
	tests := []struct {
		name    string
		target  string
		wantBuf string
		wantErr bool
	}{
		{"1", tmp + "/testdata/chart1.tgz", "alpine:3.3\n", false},
		{"2", tmp + "/testdata/chart2.tgz", "beta.opensuse.com/alpha/opensuse:42.3\n", false},
		{"3", tmp + "/testdata/.tgz", "", true},
		{"4", tmp + "/testdata/chart3.tgz", "", false},
		{"5", tmp + "/testdata/chart4.tgz", "/alpha/opensuse:42.3\n", false},
		{"6", tmp + "/testdata/chart5.tgz", "", true},
		{"7", tmp + "/testdata/chart6", "", true},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		i := &ImagesService{
			target:    "",
			formatter: fakeFormatter,
			logger:    fakeLogger,
			buffer:    buf,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := i.processTarget(tt.target); (err != nil) != tt.wantErr {
				t.Errorf("ImagesService.processTgz() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := i.buffer.String()
			if tt.wantBuf != got {
				t.Errorf("ImagesService.processTgz() buffer = %v, wantBuf %v", got, tt.wantBuf)
			}
		})
	}
	tearDownTmp()
}

func Test_sanitizeImageString(t *testing.T) {
	prepareTmp()
	tests := []struct {
		name  string
		image string
		want  string
	}{
		{"1", "\"      image:    testimage-test/image:version\"", "testimage-test/image:version"},
		{"2", "\"-      image:    testimage-test/image:version\"", "testimage-test/image:version"},
		{"3", "\" - testimage-test/image:version\"", "testimage-test/image:version"},
		{"4", "\" - {{ path }}/{{ image }}:{{ version }}\"", "{{ path }}/{{ image }}:{{ version }}"},
		{"5", "\"      image:    {{ path }}/{{ image }}:{{ version }}\"", "{{ path }}/{{ image }}:{{ version }}"},
		{"6", "\"-      image:    {{ path }}/{{ image }}:{{ version }}\"", "{{ path }}/{{ image }}:{{ version }}"},
		{"7", "testimage-test/image:version", "testimage-test/image:version"},
		{"8", "\"      image:    testhostname/testimage-test/image:version\"", "testhostname/testimage-test/image:version"},
		{"9", "\"-      image:    testhostname/testimage-test/image:version\"", "testhostname/testimage-test/image:version"},
		{"10", "\" - testhostname/testimage-test/image:version\"", "testhostname/testimage-test/image:version"},
		{"11", "\" - {{ hostname }}/{{ path }}/{{ image }}:{{ version }}\"", "{{ hostname }}/{{ path }}/{{ image }}:{{ version }}"},
		{"12", "\"      image:    {{ hostname }}/{{ path }}/{{ image }}:{{ version }}\"", "{{ hostname }}/{{ path }}/{{ image }}:{{ version }}"},
		{"13", "\"-      image:    {{ hostname }}/{{ path }}/{{ image }}:{{ version }}\"", "{{ hostname }}/{{ path }}/{{ image }}:{{ version }}"},
		{"14", "testhostname/testimage-test/image:version", "testhostname/testimage-test/image:version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeImageString(tt.image); got != tt.want {
				t.Errorf("cleanupImageString() = %v, want %v", got, tt.want)
			}
		})
	}
	tearDownTmp()
}

func prepareTmp() {
	os.MkdirAll(tmp+"/processfolder", 0777)
	os.MkdirAll(tmp+"/processfoldererror", 0777)
	os.MkdirAll(tmp+"/get", 0777)
	cpCmd := exec.Command("cp", "-R", "testdata", tmp)
	cpCmd.Run()
	tarCmd := exec.Command("tar", "zcvf", tmp+"/testdata/chart1.tgz", "--directory="+tmp+"/testdata", "chart1")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", tmp+"/testdata/chart2.tgz", "--directory="+tmp+"/testdata", "chart2")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", tmp+"/testdata/chart3.tgz", "--directory="+tmp+"/testdata", "chart3")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", tmp+"/testdata/chart4.tgz", "--directory="+tmp+"/testdata", "chart4")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", tmp+"/testdata/chart5.tgz", "--directory="+tmp+"/testdata", "chart5")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", tmp+"/processfoldererror/chart6.tgz", "--directory="+tmp+"/testdata", "chart6")
	tarCmd.Run()
	cpCmd = exec.Command("cp", tmp+"/testdata/chart1.tgz", tmp+"/processfolder")
	cpCmd.Run()
	cpCmd = exec.Command("cp", tmp+"/testdata/chart2.tgz", tmp+"/processfolder")
	cpCmd.Run()
	cpCmd = exec.Command("cp", tmp+"/testdata/chart1.tgz", tmp+"/processfoldererror")
	cpCmd.Run()
}

func tearDownTmp() {
	os.RemoveAll(tmp)
}
