package fixtures

import (
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// StartHTTPServer start http server for tests
func StartHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":1793"}
	http.HandleFunc("/alive", aliveTest)
	http.HandleFunc("/index.yaml", indexFile)
	http.HandleFunc("/chart1-2.11.0.tgz", chartTgz)
	http.HandleFunc("/chart2-1.0.1.tgz", chartTgz)
	http.HandleFunc("/chart2-0.0.0-rc1.tgz", chartTgz)
	http.HandleFunc("/chart3-0.0.1-rc1.tgz", chartTgz)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	return srv
}

// WaitForServer waits until the server is up
func WaitForServer(url string) error {
	var retry = 10
	for i := 0; i < retry; i++ {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
	}
	return errors.New("No server available")
}

func aliveTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusOK)
}

func indexFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "binary/octet-stream")
	w.Write([]byte(IndexYaml))
}

func chartTgz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "binary/octet-stream")
	w.Write(chartTGZ)
}

var chartTGZ = []byte{31, 139, 8, 0, 224, 223, 181, 91, 0, 3, 237, 193, 1, 13, 0, 0, 0, 194,
	160, 247, 79, 109, 14, 55, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 55, 3, 154, 222, 29, 39, 0, 40, 0, 0}

//Expectedcharts How many charts are in the test file
var Expectedcharts = 5

//IndexYaml test index file
var IndexYaml = `apiVersion: v1
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
    created: 2018-10-20T00:00:00.000000000Z
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
    - http://127.0.0.1:1793/chart2-0.0.0-rc1.tgz
    version: 0.0.0-rc1
  chart3:
  - apiVersion: v1
    created: 2018-12-18T00:00:00.000000000Z
    description: A Helm chart that does exist
    digest: f9a848106870c7fc8f4488f6faf178c2aa8f8f9f87fafac89819c886c86e862f
    name: chart3
    urls:
    - http://127.0.0.1:1793/chart3-0.0.1-rc1.tgz
    version: 0.0.1-rc1
  chart4:
  - apiVersion: v1
    created: 2018-12-18T00:00:00.000000000Z
    description: A Helm chart that does not exist
    digest: f9a848106870c7fc8f4488f6faf178c2aa8f8f9f87fafac89819c886c86e862f
    name: chart3
    urls:
    - http://127.0.0.1:1793/chart4-0.0.1.tgz
    version: 0.0.1-rc1
`
