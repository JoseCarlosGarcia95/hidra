# Hidra

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/hidra)](https://artifacthub.io/packages/search?repo=hidra)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidracloud/hidra)](https://goreportcard.com/report/github.com/hidracloud/hidra) [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5722/badge)](https://bestpractices.coreinfrastructure.org/projects/5722)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=bugs)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=hidracloud_hidra&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=hidracloud_hidra)

Hidra allows you to monitor the status of your services without headaches.

## Installation

### [Website monitoring with Hidra: A step-by-step guide](https://github.com/hidracloud/hidra/wiki/Website-Monitoring-with-Hidra:-A-Step%E2%80%90by%E2%80%90Step-Guide)

### Precompiled binaries

Precompiled binaries for released versions are available in the release section on GitHub. Using the latest production release binary is the recommended way of installing Hidra. You can find latest release [here](https://github.com/hidracloud/hidra/releases/latest)

### Docker images

Docker images are available on [GitHub Container Registry](https://github.com/hidracloud/hidra/pkgs/container/hidra).

### Package repositories
If you want to install Hidra easily, please use the package repositories. 

```bash
# Debian/Ubuntu
curl https://repo.hidra.cloud/apt/gpg.key | sudo apt-key add -
echo "deb [trusted=yes] https://repo.hidra.cloud/apt /" | sudo tee /etc/apt/sources.list.d/hidra.list
apt update
apt install -y hidra

# RedHat/CentOS
curl https://repo.hidra.cloud/rpm/gpg.key | sudo rpm --import -
echo "[hidra]" | sudo tee /etc/yum.repos.d/hidra.repo
echo "name=Hidra" | sudo tee -a /etc/yum.repos.d/hidra.repo
echo "baseurl=https://repo.hidra.cloud/rpm/" | sudo tee -a /etc/yum.repos.d/hidra.repo
echo "enabled=1" | sudo tee -a /etc/yum.repos.d/hidra.repo
echo "gpgcheck=1" | sudo tee -a /etc/yum.repos.d/hidra.repo

yum install -y hidra
```

After installing Hidra, test the installation by running:

```bash
hidra version
```

By default, Hidra will install one systemd unit for running hidra exporter, disabled by default. You can enable it by running:

```bash
systemctl enable hidra_exporter --now
```
You can modify the behaviour of exporter by editing the file `/etc/hidra_exporter/config.yml`, and you can add samples directly to /etc/hidra_exporter/samples/ folder, and after adding them, you can reload the service by running:

```bash
systemctl reload hidra_exporter
```


### Using install script

You can use the install script to install Hidra on your system. The script will download the latest release binary and install it in your system. You can find the script [here](https://raw.githubusercontent.com/hidracloud/hidra/main/install.sh).

```bash
sudo bash -c "$(curl -fsSL https://raw.githubusercontent.com/hidracloud/hidra/main/install.sh)"
```

### Build from source

To build Hidra from source code, you need:

- Go version 1.19 or greater.
- [Goreleaser](https://goreleaser.com)

To build Hidra, run the following command:

```bash
goreleaser release --snapshot --rm-dist
```

You can find the binaries in the `dist` folder.

## Usage

### Exporter

Hidra has support for exposing metrics to Prometheus. If you want to use Hidra in exporter mode, run:

```bash
hidra exporter /etc/hidra/exporter.yml
```

You can find an example of the configuration file [here](https://github.com/hidracloud/hidra/blob/main/configs/hidra/exporter.yml)

#### Grafana

You can find a Grafana dashboard [here](https://github.com/hidracloud/hidra/blob/main/configs/grafana)

### Test mode

Hidra has support for running in test mode. Test mode will allow you to run one time a set of samples, and check if the results are as expected. If you want to use Hidra in test mode, run:

```bash
hidra test sample1.yml sample2.yml ... samplen.yml
```

If you want to exit on error, just add the flag `--exit-on-error`.

### Sample examples

You can find some sample examples [here](https://github.com/hidracloud/hidra/blob/main/configs/hidra/samples/)

### Compose

You can find an example of a Compose file [here](https://github.com/hidracloud/hidra/blob/main/compose.yml)

## Samples

Samples are the way Hidra knows what to do. A sample is a YAML file that contains the information needed to run a test. You can find some sample examples [here](https://github.com/hidracloud/hidra/blob/main/configs/hidra/samples). You can also find a sample example below:

```yaml
# Description of the sample
description: "This is a sample to test the HTTP plugin"
# Tags is a key-value list of tags that will be added to the sample. You can add here whatever you want.
tags:
  tenant: "hidra"
# Interval is the time between each execution of the sample.
interval: "1m"
# Timeout is the time that Hidra will wait for the sample to finish.
timeout: "10s"
# Steps is a list of steps that will be executed in order.
steps:
    # Plugin is the name of the plugin that will be used to execute the step.
  - plugin: http
    # Action is the action that will be executed by the plugin.
    action: request
    # Parameters is a key-value list of parameters that will be passed to the plugin.
    parameters:
      url: https://google.com/
  - plugin: http
    action: statusCodeShouldBe
    parameters:
      statusCode: 301
```

You can find more information about plugins in next section.

## Plugins

- [browser](https://github.com/hidracloud/hidra/blob/main/docs/plugins/browser/README.md)
- [dns](https://github.com/hidracloud/hidra/blob/main/docs/plugins/dns/README.md)
- [ftp](https://github.com/hidracloud/hidra/blob/main/docs/plugins/ftp/README.md)
- [http](https://github.com/hidracloud/hidra/blob/main/docs/plugins/http/README.md)
- [icmp](https://github.com/hidracloud/hidra/blob/main/docs/plugins/icmp/README.md)
- [tcp](https://github.com/hidracloud/hidra/blob/main/docs/plugins/tcp/README.md)
- [tls](https://github.com/hidracloud/hidra/blob/main/docs/plugins/tls/README.md)
- [udp](https://github.com/hidracloud/hidra/blob/main/docs/plugins/udp/README.md)
- [string](https://github.com/hidracloud/hidra/blob/main/docs/plugins/string/README.md)

## Development

### Directory structure

The project follows the [_de facto_ standard Go project layout](https://github.com/golang-standards/project-layout) with the additions below:

- `Containerfile`, `compose.yml`, `Makefile`, `.dockerignore` and `.env.example` contain the configuration and manifests that define the development and runtime environments with [OCI](https://opencontainers.org) containers and [Compose](https://docs.docker.com/compose).
- `.github` holds the [GitHub Actions](https://github.com/features/actions) CI/CD pipelines.

### Getting started

This project comes with a containerized environment that has everything necessary to work on any platform without having to install dependencies on the developers' machines.

**TL;TR**

```Shell
make
```

#### Requirements

Before starting using the project, make sure that the following dependencies are installed on the machine:

- [Git](https://git-scm.com).
- An [OCI runtime](https://opencontainers.org), like [Podman Desktop](https://podman.io) or [Docker Desktop](https://www.docker.com/products/docker-desktop/).
- [Compose](https://docs.docker.com/compose/install/).

It is necessary to install the latest versions before continuing. You may follow the previous links to read the installation instructions.

#### Initializing

First, initialize the project and run the environment.

```Shell
make
```

Then, download third-party dependencies.

```Shell
make deps
```

You may stop the environment by running the following command.

```Shell
make down
```

### Usage

Commands must be run inside the containerized environment by starting a shell in the main container (`make shell`).

#### Running the development server

Run the following command to start the development server:

```Shell
make run
```

> Note that Git is not available in the container, so you should use it from the host machine. It is strongly recommended to use a Git GUI (like [VS Code's](https://code.visualstudio.com/docs/editor/versioncontrol) or [Fork](https://git-fork.com)) instead of the command-line interface.

#### Running tests

To run all automated tests, use the following command.

```Shell
make test
```

#### Debugging

It is possible to debug the software with [Delve](https://github.com/go-delve/delve). To run the application in debug mode, run the command below.

```Shell
make debug
```

For more advanced scenarios, such as debugging tests, you may open a shell in the container and use the Delve CLI directly.

```Shell
make shell
dlv test --listen=:2345 --headless --api-version=2 <package>
```
