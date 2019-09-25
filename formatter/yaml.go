package formatter

import (
	"bytes"
	"log"
	"strings"

	yamlencoder "gopkg.in/yaml.v3"
)

type yaml struct {
	fileName string
	l        *log.Logger
}

func newYamlFormatter(fileName string, logger *log.Logger) Formatter {
	return &yaml{
		fileName: fileName,
		l:        logger,
	}
}

func (f *yaml) Output(b bytes.Buffer) error {
	imgs := strings.Split(b.String(), "\n")
	var im Images
	for _, i := range imgs {
		if i != "" {
			im.Names = append(im.Names, i)
		}
	}
	y, err := yamlencoder.Marshal(im)
	if err != nil {
		f.l.Printf("error: cannot encode yaml")
		return err
	}
	err = writeFile(f.fileName, y, f.l)
	if err != nil {
		return err
	}
	return nil
}
