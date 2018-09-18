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

var indexYaml = `apiVersion: v1
entries:
  chart1:
  - apiVersion: v2
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing
    digest: 8cc99f9cb669171776f7c6ec66069907579be91179f9201725fc6fc6f9ef1f29
    name: chart1
    urls:
    - http://127.0.0.1:1793/chart1-2.11.0.tgz
    version: 2.11.0
  chart2:
  - apiVersion: v1
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing too
    digest: 0c76ee9b4b78cb60fcce8c00ec0f5048cbe626fcaabe48f2f8e84b029e894f49
    name: chart2
    urls:
    - http://127.0.0.1:1793/chart2-1.0.1.tgz
    version: 1.0.1
  - apiVersion: v1
    created: 2018-09-20T00:00:00.000000000Z
    description: A Helm chart for testing too
    digest: e9a545006570b7fc5e4458f6eae178c2aa8f8e9e57eafac59869c856b86e862f
    name: chart2
    urls:
    - http://127.0.0.1:1793/chart2-0.0.0-pre.tgz
    version: 0.0.0-pre
`

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

var chartTGZ = []byte{31, 139, 8, 0, 224, 223, 181, 91, 0, 3, 237, 193, 1, 13, 0, 0, 0, 194,
	160, 247, 79, 109, 14, 55, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 55, 3, 154, 222, 29, 39, 0, 40, 0, 0}
