
# go-fetch [![Build Status][build-stat]][build] [![Go Report Card][report-card-link]][report-card] [![MIT licensed][license]](./LICENSE)

I often find myself having to package up a lot of third-party and internal software in a clean way to ship to air-gapped or restricted environments. As the list of items grew, I wanted a way to easily track each component to package while maintaining some record of the versions. This tool enabled me to do that, though it is still very much a work in progress.

## Installation

```shell
go install github.com/ethanchowell/go-fetch/cmd/go-fetch@latest
```

## Quick Start

By default, `go-fetch` will look for an `artifacts.yaml` in the current working directory, and parse that file to figure out what artifacts to fetch. You can override this to look for a different file by using the `--manifest` flag.

```shell
$ go-fetch help download

Parse a given YAML manifest for artifacts that should be downloaded.

Usage:
  go-fetch download [flags]

Flags:
      --artifactory-token string   The API token for authenticating with Artifactory. Can be set from GO_FETCH_ARTIFACTORY_TOKEN.
      --bundle                     Flag to toggle if a tar.gz is generated with the same name as the --manifest flag.
      --github-token string        The API token for authenticating with GitHub. Can be set from GO_FETCH_GITHUB_TOKEN.
      --gitlab-token string        The API token for authenticating with GitLab. Can be set from GO_FETCH_GITLAB_TOKEN.
  -h, --help                       help for download
      --manifest string            Path to the manifest containing artifacts to download. Can be set from GO_FETCH_MANIFEST. (default "./artifacts.yaml")
```

The download is driven by a set of "providers" which, currently, are the following

| Provider    | Functional |
|-------------|------------|
| Artifactory | [ ]        |
| Docker      | [x]        |
| Github      | [x]        |
| GitLab      | [ ]        |
| Helm        | [x]        |
| HTTP GET    | [ ]        |

You can mix usage of these providers in the same file.

### Docker Images

Each entry in the manifest must define one of these providers. An example manifest to pull, and save a set of docker image would look like this

```yaml
# The top-level directory to save things
target: artifacts
releases:
  - repo:
      name: docker.io
      provider: docker
    artifacts:
      - ubuntu:20.04
      - ubuntu:focal
  - tag: 14.7.0-debian-11-r10 # If the image names aren't tagged, define it here
    repo:
      name: bitnami
      provider: docker
    artifacts:
      - postgresql
```

Running `go-fetch download` against this manifest will download the artifacts into the following structure.

```shell
$ tree artifacts 
artifacts
├── bitnami
│ └── 14.7.0-debian-11-r10
│     └── postgresql.tar.gz
├── docker.io
│ ├── ubuntu-20.04.tar.gz
│ └── ubuntu-focal.tar.gz
└── sha265sum.txt

3 directories, 4 files
```

with the `sha265sum` containing the SHA256 checksum of each artifact.

Notice that you can optionally supply the tag from a key, or as part of the artifact name. This is only supported for Docker images and Helm charts.

### GitHub Releases

You can fetch Github Releases with a manifest like the following.

```yaml
# The top-level directory to save things
target: artifacts
releases:
  - tag: v3.11.2
    repo:
      name: helm/helm
      provider: github
    artifacts:
      - helm-v3.11.2-linux-amd64.tar.gz.asc
  - tag: v0.38.2
    repo:
      name: aquasecurity/trivy
      provider: github
    artifacts:
      - trivy_0.38.2_Linux-64bit.tar.gz
```

Running `go-fetch download` against this manifest will download the artifacts into the following structure.

```shell
$ tree artifacts 
artifacts
├── aquasecurity
│ └── trivy
│     └── v0.38.2
│         └── trivy_0.38.2_Linux-64bit.tar.gz
├── helm
│ └── helm
│     └── v3.11.2
│         └── helm-v3.11.2-linux-amd64.tar.gz.asc
└── sha265sum.txt

6 directories, 3 files
```

### Helm Charts

You can fetch Helm charts with a manifest like the following.

```yaml
# The top-level directory to save things
target: artifacts
releases:
  - tag: 12.2.2 # If the chart names aren't tagged, define it here
    repo:
      name: bitnami # From helm repo add bitnami https://charts.bitnami.com/bitnami 
      provider: helm
    artifacts:
      - postgresql
  - tag: 11.1.4
    repo:
      name: https://charts.bitnami.com/bitnami # Doesn't need to be from `helm repo add`
      provider: helm
    artifacts:
      - postgresql-ha
```

Running `go-fetch download` against this manifest will download the artifacts into the following structure.

```shell
$ tree artifacts 
artifacts
├── bitnami
│ └── 12.2.2
│     └── postgresql-12.2.2.tgz
├── helm # If a URL is provided, we'll write to `helm`
│ └── 11.1.4
│     └── postgresql-ha-11.1.4.tgz
└── sha265sum.txt

4 directories, 3 files
```

[build-stat]: https://github.com/ethanchowell/go-fetch/actions/workflows/go.yml/badge.svg
[build]: https://github.com/ethanchowell/go-fetch/actions/workflows/go.yml
[report-card-link]: https://goreportcard.com/badge/github.com/ethanchowell/go-fetch
[report-card]: https://goreportcard.com/report/github.com/ethanchowell/go-fetch
[license]: https://img.shields.io/badge/license-MIT-blue.svg
