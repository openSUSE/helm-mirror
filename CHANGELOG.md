# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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