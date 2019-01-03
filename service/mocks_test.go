package service

import (
	"bytes"

	"github.com/pkg/errors"
)

var errImplemented = errors.New("not implemented")

type mockFormatter struct {
}

func (m *mockFormatter) Output(buffer bytes.Buffer) error {
	if buffer.String() == "test" {
		return errors.New("not implemented")
	}
	return nil
}

var valuesYaml = `---
kube:
  external_ips: []
  storage_class:
    persistent: "persistent"
    shared: "shared"
  registry:
    hostname: "beta.opensuse.com"
    username: ""
    password: ""
  organization: "alpha"
image: "opensuse"
version: "42"
`
