# IntegreSQL Client Library for Golang

Client library for interacting with a [`IntegreSQL` server](https://github.com/allaboutapps/integresql), managing isolated PostgreSQL databases for your integration tests.

## Overview [![](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/allaboutapps/integresql-client-go?tab=doc) [![](https://goreportcard.com/badge/github.com/allaboutapps/integresql-client-go)](https://goreportcard.com/report/github.com/allaboutapps/integresql-client-go) ![](https://github.com/allaboutapps/integresql-client-go/workflows/build/badge.svg?branch=master)

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contributing](#contributing)
    - [Development setup](#development-setup)
    - [Development quickstart](#development-quickstart)
- [Maintainers](#maintainers)
- [License](#license)

## Background

See [IntegreSQL: Background](https://github.com/allaboutapps/integresql#background)

## Install

Install the `IntegreSQL` client for Go using `go get` or by simply importing the library in your testing code (go modules, see below):

```bash
go get github.com/allaboutapps/integresql-client-go
```

## Configuration

The `IntegreSQL` client library requires little configuration which can either be passed via the `ClientConfig` struct or parsed from environment variables automatically. The following settings are available:

| Description                                                | Environment variable            | Default                        | Required |
| ---------------------------------------------------------- | ------------------------------- | ------------------------------ | -------- |
| IntegreSQL: base URL of server `http://127.0.0.1:5000/api` | `INTEGRESQL_CLIENT_BASE_URL`    | `"http://integresql:5000/api"` |          |
| IntegreSQL: API version of server                          | `INTEGRESQL_CLIENT_API_VERSION` | `"v1"`                         |          |


## Usage

If you want to take a look on how we integrate IntegreSQL - 🤭 - please just try our [go-starter](https://github.com/allaboutapps/go-starter) project or take a look at our [testing setup code](https://github.com/allaboutapps/go-starter/blob/master/internal/test/testing.go). 

In general setting up the `IntegreSQL` client, initializing a PostgreSQL template (migrate + seed) and retrieving a PostgreSQL test database goes like this:

```go
package yourpkg

import (
    "github.com/allaboutapps/integresql-client-go"
    "github.com/allaboutapps/integresql-client-go/pkg/util"
)

func doStuff() error {
    c, err := integresql.DefaultClientFromEnv()
    if err != nil {
        return err
    }

    // compute a hash over all database related files in your workspace (warm template cache)
    hash, err := hash.GetTemplateHash("/app/scripts/migrations", "/app/internal/fixtures/fixtures.go")
    if err != nil {
        return err
    }

    template, err := c.InitializeTemplate(context.TODO(), hash)
    if err != nil {
        return err
    }

    // Use template database config received to initialize template
    // e.g. by applying migrations and fixtures

    if err := c.FinalizeTemplate(context.TODO(), hash); err != nil {
        return err
    }

    test, err := c.GetTestDatabase(context.TODO(), hash)
    if err != nil {
        return err
    }

    // Use test database config received to run integration tests in isolated DB
}
```

A very basic example has been added as the `cmd/cli` executable, you can build it using `make cli` and execute `integresql-cli` afterwards.

## Contributing

Pull requests are welcome. For major changes, please [open an issue](https://github.com/allaboutapps/integresql/issues/new) first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### Development setup

`IntegreSQL` requires the following local setup for development:

- [Docker CE](https://docs.docker.com/install/) (19.03 or above)
- [Docker Compose](https://docs.docker.com/compose/install/) (1.25 or above)

The project makes use of the [devcontainer functionality](https://code.visualstudio.com/docs/remote/containers) provided by [Visual Studio Code](https://code.visualstudio.com/) so no local installation of a Go compiler is required when using VSCode as an IDE.

Should you prefer to develop the `IntegreSQL` client library without the Docker setup, please ensure a working [Go](https://golang.org/dl/) (1.14 or above) environment has been configured as well as an `IntegreSQL` server and a a PostgreSQL instance are available (tested against PostgreSQL version 12 or above, but *should* be compatible to lower versions) and the appropriate environment variables have been configured as described in the [Install](#install) section.

### Development quickstart

1. Start the local docker-compose setup and open an interactive shell in the development container:

```bash
# Build the development Docker container, start it and open a shell
./docker-helper.sh --up
```

2. Initialize the project, downloading all dependencies and tools required (executed within the dev container):

```bash
# Init dependencies/tools
make init

# Build executable (generate, format, build, vet)
make
```

3. Execute project tests:

```bash
# Execute tests
make test
```

## Maintainers

- [Nick Müller - @MorpheusXAUT](https://github.com/MorpheusXAUT)
- [Mario Ranftl - @majodev](https://github.com/majodev)

## License

[MIT](LICENSE) © 2020 aaa – all about apps GmbH | Nick Müller | Mario Ranftl and the `IntegreSQL` project contributors