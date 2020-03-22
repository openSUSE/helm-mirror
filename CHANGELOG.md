# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## v0.3.2

- Support relative chart URL in `index.yaml` (such as those created by Harbor).

## v0.3.1

- Update to use go modules
- Update to use hel 2.16.1 to fix CVE-2019-18658

bsc#1156646

## v0.3.0

- New features: get latest and specific charts
  - Getting charts now only downloads the altest versions of the charts.
  - The --all-versions flags allows to download all versions of the charts.
  - The flags --chart-name and --chart-version allow the user to only get the desired chart.

## v0.2.4

- updated release steps
- updated install script

## v0.2.3

- fixes issue while getting the latest release

## v0.2.2

- fixes issue with go module when installing with `helm plugin install`

## v0.2.1

- fixes empty archive files and usage of ignore-errors flag

## v0.2.0

### Added

- `mirror inspect-images` flag `--output` usage updated, flag `--file-name` no longer needed.
  - `file=filename`
  - `json=filename.json`
  - `yaml=filename.yaml`
  - `skopeo=filename.yaml`

- `helm-mirror` has a new flag `--new-root-url` new root url of the chart repository.
  (eg: `https://mirror.local.lan/charts`). This will allow users to set the name of
  their mirror server when getting all the charts.

- `downloaded-index.yaml` file changes it's name to `index.yaml` to allow users to host quickly
  a mirror chart server.
  

## v0.1.0

### Added

- `mirror [chart-repo] [target-folder]` this command takes a chart repository and downloads all
  chart found in there and downloads them to a local target folder.

- `mirror inspect-images [target]` this command takes the target and extracts all container
  images being used in it. It can be a single chart or a folder with multiple charts. If
  some values are not present to process the chart you can use `--ignore-errors` flag to
  render the charts and get the container image anyway. (this can output inconsistent data)

- Use Helm configuration settings for Chart Repository.
  - --ca-file
  - --cert-file
  - --key-file
  - --password
  - --username

- CI steps with [Travis CI](https://travis-ci.org) and Code Coverage with [CodeCov](https://codecov.io)

[Unreleased]: https://github.com/openSUSE/helm-mirror/compare/v0.1.0...HEAD
