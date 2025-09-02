# MrAndreID / Go Application Programming Interface (API) Gateway

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The `MrAndreID/GoAPIGateway` is a skeleton uses the Go Programming Language (GoLang) with The Echo Framework for The Application Programming Interface (API) Gateway.

## Table of Contents

* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [Versioning](#versioning)
* [Authors](#authors)
* [Contributing](#contributing)
* [Official Documentation for Go Language](#official-documentation-for-go-language)
* [License](#license)

## Requirements

To use The `MrAndreID/GoAPIGateway`, you must ensure that you meet the following requirements:
- [Go](https://golang.org/) >= 1.24

## Installation

To use The `MrAndreID/GoAPIGateway`, you must follow the steps below:
- Clone a Repository
```git
# git clone https://github.com/MrAndreID/goapigateway.git
```
- Get Dependancies
```go
# go mod download
# go mod tidy
```
- Create .env file from .env.example (Linux)
```sh
# cp .env.example .env
```
- Configuring .env file

## Usage

To use The `MrAndreID/GoAPIGateway`, you must ensure that you meet the following requirements:
- Directory Structure The `MrAndreID/GoAPIGateway`
| Name                    | Description                                               |
| :---------------------- | :-------------------------------------------------------- |
| `application`           | Initialization of Echo Framework, Middleware, and Routes. |
| `configs`               | Configuration from Env File                               |
| `internal/handlers`     | HTTP Handlers                                             |
| `internal/services`     | Main Business Logic                                       |
| `internal/repositories` | Connector to Database or API External                     |
| `internal/types`        | Data                                                      |
- Run The `MrAndreID/GoAPIGateway`
```go
# go run main.go
```
- Run The `MrAndreID/GoAPIGateway` with Docker
```docker
# docker build --no-cache -t goapigateway:1.0.0 .
# docker run --name goapigateway --restart=always -d -p -v /path/to/folder:/app/storages -v /path/to/folder:/app/tests/storages 10000:10000 goapigateway:1.0.0
```
- Set The `MrAndreID/GoAPIGateway` to Maintenance Mode in Storages Folder
```sh
# touch storages/maintenance.flag
```

## Versioning

I use [Semanting Versioning](https://semver.org/). For the versions available, see the tags on this repository. 

## Authors

- **Andrea Adam** - [MrAndreID](https://github.com/MrAndreID)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
Please make sure to update tests as appropriate.

## Official Documentation for Go Language

Documentation for Go Language can be found on the [Go Package website](https://pkg.go.dev/).

## License

The `MrAndreID/GoAPIGateway` is released under the [MIT License](https://opensource.org/licenses/MIT). See the `LICENSE` file for more information.
