package service

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"testing"

	"github.com/openSUSE/helm-mirror/formatter"
)

var (
	buff          bytes.Buffer
	fakeFormatter = &mockFormatter{}
)

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
	dir, err := prepareTmp()
	if err != nil {
		t.Errorf("loading testdata: %s", err)
	}
	defer os.RemoveAll(dir)
	processPath := path.Join(dir, "processfolder")
	errorPath := path.Join(dir, "processfoldererror")
	testdataPath := path.Join(dir, "testdata")
	processTgzPath := path.Join(dir, "processtgz")
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
		{"1", fields{processPath, buff, false}, false},
		{"2", fields{path.Join(testdataPath, "chart6"), buff, false}, true},
		{"3.1", fields{path.Join(testdataPath, "chart1"), buff, false}, false},
		{"3.2", fields{path.Join(processTgzPath, "chart1.tgz"), buff, false}, false},
		{"4", fields{path.Join(dir, "mr", "mzxyptlk"), buff, false}, true},
		{"5", fields{errorPath, buff, true}, false},
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
	dir, err := prepareTmp()
	if err != nil {
		t.Errorf("loading testdata: %s", err)
	}
	defer os.RemoveAll(dir)
	processPath := path.Join(dir, "processfolder")
	processTgzPath := path.Join(dir, "processtgz")
	type fields struct {
		target    string
		imageFile string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"1", fields{processPath, "images"}, false},
		{"2", fields{processTgzPath, "images"}, true},
		{"3", fields{path.Join(processTgzPath, "chart1.tgz"), "images"}, true},
		{"4", fields{path.Join(processPath, "chart6.tgz"), "images"}, true},
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
}

func TestImagesService_processTarget(t *testing.T) {
	dir, err := prepareTmp()
	if err != nil {
		t.Errorf("loading testdata: %s", err)
	}
	defer os.RemoveAll(dir)
	processTgzPath := path.Join(dir, "processtgz")
	tests := []struct {
		name    string
		target  string
		verbose bool
		wantBuf string
		wantErr bool
	}{
		{"1", path.Join(processTgzPath, "chart1.tgz"), false, "alpine:3.3\n", false},
		{"2", path.Join(processTgzPath, "chart2.tgz"), false, "beta.opensuse.com/alpha/opensuse:42.3\n", false},
		{"3", path.Join(processTgzPath, ".tgz"), false, "", true},
		{"4", path.Join(processTgzPath, "chart3.tgz"), false, "", false},
		{"5", path.Join(processTgzPath, "chart4.tgz"), false, "/alpha/opensuse:42.3\n", false},
		{"6", path.Join(processTgzPath, "chart5.tgz"), false, "", true},
		{"7", path.Join(processTgzPath, "chart6"), false, "", true},
		{"8", path.Join(processTgzPath, "chart6"), true, "", true},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		i := &ImagesService{
			target:    "",
			formatter: fakeFormatter,
			logger:    fakeLogger,
			buffer:    buf,
			verbose:   tt.verbose,
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
}

func Test_sanitizeImageString(t *testing.T) {
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
}

func prepareTmp() (string, error) {
	dir, err := ioutil.TempDir("", "helmmirror")
	if err != nil {
		return "", err
	}
	testdataPath := path.Join(dir, "testdata")
	processPath := path.Join(dir, "processfolder")
	processTgzPath := path.Join(dir, "processtgz")
	errorPath := path.Join(dir, "processfoldererror")
	getPath := path.Join(dir, "get")
	os.MkdirAll(processPath, 0777)
	os.MkdirAll(processTgzPath, 0777)
	os.MkdirAll(errorPath, 0777)
	os.MkdirAll(getPath, 0777)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	err = os.Symlink(path.Join(wd, "testdata"), testdataPath)
	if err != nil {
		return "", err
	}
	tarCmd := exec.Command("tar", "zcvf", path.Join(processTgzPath, "chart1.tgz"), "--directory="+testdataPath, "chart1")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", path.Join(processTgzPath, "chart2.tgz"), "--directory="+testdataPath, "chart2")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", path.Join(processTgzPath, "chart3.tgz"), "--directory="+testdataPath, "chart3")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", path.Join(processTgzPath, "chart4.tgz"), "--directory="+testdataPath, "chart4")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", path.Join(processTgzPath, "chart5.tgz"), "--directory="+testdataPath, "chart5")
	tarCmd.Run()
	tarCmd = exec.Command("tar", "zcvf", path.Join(errorPath, "chart6.tgz"), "--directory="+testdataPath, "chart6")
	tarCmd.Run()

	err = os.Symlink(path.Join(processTgzPath, "chart1.tgz"), path.Join(processPath, "chart1.tgz"))
	if err != nil {
		return "", err
	}
	err = os.Symlink(path.Join(processTgzPath, "chart2.tgz"), path.Join(processPath, "chart2.tgz"))
	if err != nil {
		return "", err
	}
	err = os.Symlink(path.Join(processTgzPath, "chart1.tgz"), path.Join(errorPath, "chart1.tgz"))
	if err != nil {
		return "", err
	}
	return dir, err
}
