package formatter

import (
	"bytes"
	"log"
	"strings"

	"github.com/containers/image/types"
	"github.com/docker/distribution/reference"
	yamlencoder "gopkg.in/yaml.v3"
)

type skopeo struct {
	fileName string
	l        *log.Logger
}

func newSkopeoFormatter(fileName string, logger *log.Logger) Formatter {
	return &skopeo{
		fileName: fileName,
		l:        logger,
	}
}

func (f *skopeo) Output(b bytes.Buffer) error {
	images := strings.Split(b.String(), "\n")
	var registries Registries
	registries = make(map[string]Registry)
	for _, i := range images {
		if i != "" {
			ref, err := reference.ParseNormalizedNamed(i)
			if err != nil {
				f.l.Printf("error: parsing image %s", i)
				continue
			}
			registry := reference.Domain(ref)
			image := reference.Path(ref)
			tag := ""
			if named, ok := ref.(reference.NamedTagged); ok {
				tag = named.Tag()
			}
			if r, ok := registries[registry]; ok {
				if im, ok := r.Images[image]; ok {
					r.Images[image] = append(im, tag)
				} else {
					r.Images[image] = []string{tag}
				}
			} else {
				registries[registry] = Registry{
					Images:      make(map[string][]string),
					Credentials: types.DockerAuthConfig{},
					CertDir:     "",
					TLSVerify:   false,
				}
				registries[registry].Images[image] = []string{tag}
			}
		}
	}
	y, err := yamlencoder.Marshal(registries)
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

//TODO use Skopeo's source code structures

// Registry definition of a registry to be used by Skopeo
type Registry struct {
	Images      map[string][]string    `yaml:"images"`
	Credentials types.DockerAuthConfig `yaml:"credentials,omitempty"`
	TLSVerify   bool                   `yaml:"tls-verify,omitempty"`
	CertDir     string                 `yaml:"cert-dir,omitempty"`
}

// Registries defines a map of Registries
type Registries map[string]Registry
