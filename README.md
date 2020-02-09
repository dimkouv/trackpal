# Trackpal

A location tracking application


## Quickstart

This application consists of a RESTful API that is capable
of logging a stream of coordinates and then exposing them
in order to track the things you love (vehicles, people, ...).


## Architecture

Initially we're starting with a simple approach. The server is
written in Go and Postgres will be our database.


## Tools & Dependencies

golangci-lint

```bash
if ! [ -x "$(command -v golangci-lint)" ]; then
  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin v1.21.0
fi
```


docsify

```bash
npm i docsify-cli -g
```


postgres

```bash
docker run --name postgres-local -e POSTGRES_PASSWORD=masterkey -d postgres
```
