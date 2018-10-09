package formatter

import (
	"bytes"
	"log"
	"os"
)

type stdout struct {
	l *log.Logger
}

func newStdoutFormatter(logger *log.Logger) Formatter {
	return &stdout{
		l: logger,
	}
}

func (s *stdout) Output(b bytes.Buffer) error {
	_, err := b.WriteTo(os.Stdout)
	if err != nil {
		s.l.Printf("error: cannot write to stdout: %s", err)
		return err
	}
	return nil
}
