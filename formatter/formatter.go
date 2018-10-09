package formatter

import (
	"bytes"
	"io/ioutil"
	"log"
)

//Formatter defines the behavior for a Formatter
type Formatter interface {
	Output(buffer bytes.Buffer) error
}

//Type definition of formatter type enum
type Type int

//Enum for Formatter
const (
	StdoutType Type = 1 << iota
	FileType
	JSONType
	YamlType
)

//NewFormatter returns a new instance of formatter
func NewFormatter(t Type, fileName string, logger *log.Logger) Formatter {
	switch t {
	case StdoutType:
		return newStdoutFormatter(logger)
	case FileType:
		return newFileFormatter(fileName, logger)
	case JSONType:
		return newJSONFormatter(fileName, logger)
	case YamlType:
		return newYamlFormatter(fileName, logger)
	default:
		return newStdoutFormatter(logger)
	}
}

func writeFile(name string, content []byte, log *log.Logger) error {
	err := ioutil.WriteFile(name, content, 0666)
	if err != nil {
		log.Printf("cannot write files %s: %s", name, err)
		return err
	}
	return nil
}

//Images struct for YAML and JSON output
type Images struct {
	Names []string `json,yaml:"names,omitempty"`
}
