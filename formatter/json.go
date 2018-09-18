package formatter

import (
	"bytes"
	jsonencoding "encoding/json"
	"log"
	"strings"
)

type json struct {
	fileName string
	l        *log.Logger
}

func newJSONFormatter(fileName string, logger *log.Logger) Formatter {
	return &json{
		fileName: fileName,
		l:        logger,
	}
}

func (f *json) Output(b bytes.Buffer) error {
	imgs := strings.Split(b.String(), "\n")
	var im Images
	for _, i := range imgs {
		if i != "" {
			im.Names = append(im.Names, i)
		}
	}
	j, err := jsonencoding.Marshal(im)
	if err != nil {
		f.l.Printf("error: cannot encode json")
		return err
	}
	err = writeFile(f.fileName, j, f.l)
	if err != nil {
		return err
	}
	return nil
}
