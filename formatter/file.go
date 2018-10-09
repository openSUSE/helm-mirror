package formatter

import (
	"bytes"
	"log"
)

type file struct {
	fileName string
	l        *log.Logger
}

func newFileFormatter(fileName string, logger *log.Logger) Formatter {
	return &file{
		fileName: fileName,
		l:        logger,
	}
}

func (f *file) Output(b bytes.Buffer) error {
	err := writeFile(f.fileName, b.Bytes(), f.l)
	if err != nil {
		return err
	}
	return nil
}
