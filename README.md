# Trackpal
![Go](https://github.com/dimkouv/trackpal/workflows/Go/badge.svg)

A location tracking application.


## Quickstart

This application consists of a RESTful API that is capable
of logging a stream of coordinates and then exposing them in 
order to track and get alerted when the things you love (vehicles, people, ...)
are moving while they're supposed not to. Client implementations are not ready yet.

Check `make help` for available commands.

## Architecture

Initially we're starting with a simple approach. The server is
written in Go, Postgres will serve as our database and http will be handled by a `net/http` server .

We have an approach with 3 layers (server, service, repository). The repositories
are responsible for storage operations, services are using the repositories
and contain all business logic, server uses services to serve http requests.

This approach helps us easily update the storage (e.g. replace postgres with mysql <*not*>),
allows to create mock storage implementation for fast unit tests. The server is separated from the
business logic in order to make it easy for adapting to new standards, for example replace net/http
 with lambdas or grpc.


```text
internal/server  -->  internal/services  -->  internal/repository
                      |                       |
                      | model_service.go      | model.go          (iface)   
                                              | model_mock.go     (mock impl)
                                              | model_postgres.go (pg impl)
                                              | model_redis.go    (redis impl)


- inits services      - returns enum errors   - returns repo errors
- http handlers       - returns bytes         - returns models
- routes              - repo errors logging   - no logging
```


## Tools & Dependencies

### golangci-lint

Install golangci-lint in order to be able to run `make lint`

```bash
if ! [ -x "$(command -v golangci-lint)" ]; then
  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin v1.21.0
fi
```


### postgres

You can run postgres locally with docker.

```bash
docker run --name postgres-local \
    -e POSTGRES_USER=master \
    -e POSTGRES_DB=trackpal \
    -e POSTGRES_PASSWORD=masterkey \
    -p 5432:5432 -d postgres
```

After restarting your system, you can start 
the instance with `docker start postgres-local`
