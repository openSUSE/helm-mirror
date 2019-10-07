% helm-mirror(1) # helm-mirror - Mirror chart repositories
% SUSE LLC
% OCTOBER 2018
# NAME
helm-mirror - helm-mirror chart repositories

# SYNOPSIS
**helm-mirror**
[**--help**|**-h**]
[**version**]
[**inspect-images**]
[**--ca-file**]
[**--cert-file**]
[**--chart-name**]
[**--chart-version**]
[**--ignore-errors**]
[**--key-file**]
[**--new-root-url**]
[**--password**]
[**--username**]
[**--verbose**|**-v**]
*command* [*args*]

# DESCRIPTION
**helm-mirror** is a [Helm][1] plugin that allows the mirroring of a Chart
repository, with the index file and all referenced charts.

**helm-mirror** will allow users to pass through the same options
as [Helm][1] does to connect to the repository and download the charts.

**helm-mirror** will also allow users to inspect the charts and extract from them the
container images that they use. This will be allowed by the sub-command
**helm-mirror-inspect-images**(1).

The index file is a yaml that contains a list of charts in this format.
Example:

```
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

-https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz

-https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

into your destination folder.

# GLOBAL OPTIONS

**-h, --help**
  Print usage statement.

**-v, --verbose**
  Verbose output

**--ca-file**
  Verify certificates of HTTPS-enabled servers using this CA bundle

**--cert-file**
  Identify HTTPS client using this SSL certificate file

**--chart-name**
  Name of the desired chart to download

**--chart-version**
  Version of the desired chart to download, needs the `--chart-name` option

**-i, --ignore-errors**
  Ignores errors while downloading or processing charts

**--key-file**
  Identify HTTPS client using this SSL key file

**--new-root-url**
  New root url of the chart repository (eg: `https://mirror.local.lan/charts`)

**--password**
  Chart repository password

**--username**
  Chart repository username

# COMMANDS

**inspect-images**
  Extract the images from the a target. See **helm-mirror-inspect-images**(1) for more detailed usage
  information.

**version**
  Print current version of software. See **helm-mirror-version**(1) for more detailed
  usage information.

**help**
  Print usage statements. See **helm-mirror-help**(1)
  for more detailed usage information.

# EXAMPLES
This will read the given chart repository and download all the Charts that are found into the
given target folder. target folder has to be always an absolute path.

`% helm-mirror https://yourorg.com/charts /yourorg/charts`

This will download all the charts and all the available versions
of the charts.

`% helm-mirror https://yourorg.com/charts /yourorg/charts --all-charts`

This will download the latest version of the chart `nginx`.

`% helm-mirror https://yourorg.com/charts /yourorg/charts --chart-name nginx`

This will download the version `2.14.3` of the chart `nginx`.

`% helm-mirror https://yourorg.com/charts /yourorg/charts --chart-name nginx --chart-version 2.14.3`


# SEE ALSO
**helm-mirror-inspect-images**(1),
**helm-mirror-help**(1),
**helm-mirror-version**(1)

[1]: https://docs.helm.sh
