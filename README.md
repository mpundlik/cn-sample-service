# CN sample service

[![Build Status](https://travis-ci.org/ligato/cn-sample-service.svg?branch=master)](https://travis-ci.org/ligato/cn-sample-service)
[![Coverage Status](https://coveralls.io/repos/github/ligato/cn-sample-service/badge.svg?branch=master)](https://coveralls.io/github/ligato/cn-sample-service?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ligato/cn-sample-service)](https://goreportcard.com/report/github.com/ligato/cn-sample-service)
[![GoDoc](https://godoc.org/github.com/ligato/cn-sample-service?status.svg)](https://godoc.org/github.com/ligato/cn-sample-service)
[![GitHub license](https://img.shields.io/badge/license-Apache%20license%202.0-blue.svg)](https://github.com/ligato/cn-sample-service/blob/master/LICENSE)

The CN sample service showcases the extensibility of [cn-infra](https://github.com/ligato/cn-infra). The cn-infra repository has its own examples
but this dedicated repository demonstrates the [dependency management](Glide) and building using [makefiles](Makefile).
Use this repository as skeleton for your software projects (copy&paste it at the very beginning).

![sample_service](docs/imgs/sample_service.png "Sample service plugins")


The sample service repository contains:
* [Hello World!](cmd/helloworld) - minimalistic extension of the cn-infra
* [Cassandra](cmd/cassandra) - sample service which makes REST API calls to interact with Cassandra
* Flavors/Plugins are used from [cn-infra](https://github.com/ligato/cn-infra)
* [Core](https://github.com/ligato/cn-infra/tree/master/core) from cn-infra - lifecycle management of plugins (loading, 
initialization, unloading)

## Quickstart

The basic steps to setup a project:

Initialize a tool for management dependencies. This example uses [glide](https://github.com/Masterminds/glide).
We assume that glide is already installed if not, follow the instruction in its README.

```
glide init
```

Modify the content of the `glide.yaml` that defines dependencies of the project.
It is recommanded to pin dependencies to a particular commit id or a tag. The initial content
might look like this:

```yaml
package: github.com/ligato/cn-sample-service
import:
- package: github.com/ligato/cn-infra
  version: 7657681a90ee7630248e8bf312dbfae92d82635e
```

Download initial set of dependencies

```
glide install --strip-vendor
```

Once the initial set of dependencies is downloaded we can move to the writing of a custom plugin.
The common tasks related to project development such as building, updating of dependencies, running of static
analysis and so on... can be automated using Makefile. Take a look at the [Makefile](Makefile) in this repository.

Examples of custom plugins:
 - [HelloWorld](cmd/helloworld) - minimalistic sample to get started writing your own plugin.
 - [Cassandra](cmd/cassandra) - detailed implementation of REST API based micro-service which interacts with Cassandra via cn-infra SQL-like API layer.

## Makefile
Source codes in this repository have own Makefile. This Makefile can be modified and extended based on requirements
of a particular software project.

## Glide
Glide.yaml contains import of the cn-infra therefore vendor directory in source codes of this repository
will contain all transitive dependencies of the cn-infra. Note, if you use just subset of these vendor packages
golang will statically build/link only the subset (not all packages in vendor).