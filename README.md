# Helm Mirror Plugin

Helm plugin used to mirror repositories

## Usage

Mirror Helm Charts from an index file into a local folder.

For example:

`helm mirror https://yourorg.com/charts /yourorg/charts`

This will download the index file and the charts into
the folder indicated.

The index file is a yaml that contains a list of
charts in this format. Example:

```yaml
apiVersion: v1
entries:
  chart:
  - apiVersion: 1.0.0
    created: 2018-08-08T00:00:00.00000000Z
    description: A Helm chart for your application
    digest: 3aa68d6cb66c14c1fcffc6dc6d0ad8a65b90b90c10f9f04125dc6fcaf8ef1b20
    name: chart
    urls:
    - https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz
  chart2:
  - apiVersion: 1.0.0
    created: 2018-08-08T00:00:00.00000000Z
    description: A Helm chart for your application
    digest: 7ae62d60b61c14c1fcffc6dc670e72e62b91b91c10f9f04125dc67cef2ef0b21
    name: chart
    urls:
    - https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz
```

This will download these charts:

- https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz

- https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

into your destination folder.

Usage:

```
  helm mirror [Repo URL] [Destination Folder] [flags]
  helm mirror [command]
```

Available Commands:
  help           Help about any command
  inspect-images Extract all the images of the Helm Charts.
  version        Show version of the helm mirror plugin

Flags:

```
      --ca-file string     verify certificates of HTTPS-enabled servers using this CA bundle
      --cert-file string   identify HTTPS client using this SSL certificate file
  -h, --help               help for mirror
  -i, --ignore-errors      ignores errors whiles processing charts. (Exit Code: 2)
      --key-file string    identify HTTPS client using this SSL key file
      --password string    chart repository password
      --username string    chart repository username
  -v, --verbose            verbose output
```

Use `helm mirror [command] --help` for more information about a command.

## Commands

### inspect-images

Extract all the images of the Helm Chart or
the Helm Charts in the folder provided. This command dumps
the images on `stdout` by default, for more options check
`output flag`. Example:

- helm mirror inspect-images /tmp/helm

- helm mirror inspect-images /tmp/helm/app.tgz

The [folder|tgzfile] has to be a full path.

#### Usage

`mirror inspect-images [folder|tgzfile] [flags]`

#### Flags

  -h, --help               help for inspect-images

  --file-name string   set the name of the output file. (default "images.out")

```shell
helm mirror inspect-images /tmp/helm -o file --file-name images.txt
helm mirror inspect-images /tmp/helm -o json --file-name images.json
helm mirror inspect-images /tmp/helm -o yaml --file-name images.yaml
```

  -o, --output string      choose an output for the list of images.(default "stdout")

- stdout: prints all images on stdout
- file: outputs all images to a file. (View file-name flag)
- json: outputs all images to a file in JSON format. (View file-name flag)
- yaml: outputs all images to a file in YAML format. (View file-name flag)

```shell
helm mirror inspect-images /tmp/helm --output stdout
helm mirror inspect-images /tmp/helm -o stdout
helm mirror inspect-images /tmp/helm -o file
helm mirror inspect-images /tmp/helm -o json
helm mirror inspect-images /tmp/helm -o yaml
```

#### Global Flags

  -i, --ignore-errors   ignores errors whiles processing charts. (Exit Code: 2)
  -v, --verbose         verbose output

### version

Displays the current version of mirror.

## Install

Using Helm plugin manager (> 2.3.x)

`helm plugin install https://github.com/openSUSE/helm-mirror --version master`

## Test

Clone repository into your $GOPATH. You can also use go get:

`go get github.com/openSUSE/helm-mirror`

### Prerequisites

- Have [GO](https://golang.org/) installed.

- We use [Go Dep](https://github.com/golang/dep) as dependency manager so you'll need that too.

### Bootstrap

Get all dependencies running:

`dep ensure`

### Runing tests

To run test on this package simply run:

`make test`

#### Testign with Docker

`make test.unit`

## Building

Be sure you have all prerequisites, then build the binary by simply running

`make mirror`

your binary will be stored under `bin` folder