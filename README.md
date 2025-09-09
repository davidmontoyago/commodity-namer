# commodity-namer

[![Develop](https://github.com/davidmontoyago/commodity-namer/actions/workflows/develop.yaml/badge.svg)](https://github.com/davidmontoyago/commodity-namer/actions/workflows/develop.yaml) [![Go Coverage](https://raw.githubusercontent.com/wiki/davidmontoyago/commodity-namer/coverage.svg)](https://raw.githack.com/wiki/davidmontoyago/commodity-namer/coverage.html) [![Go Reference](https://pkg.go.dev/badge/github.com/davidmontoyago/commodity-namer.svg)](https://pkg.go.dev/github.com/davidmontoyago/commodity-namer)

Consistent structured naming for infra resources.

### Getting started

```go
import "github.com/davidmontoyago/commodity-namer"

type MyInfra struct {
  // Embed the namer
  Namer
}

func NewMyInfra() *MyInfra {
  return &MyInfra{
    // Set the base name
    Namer: namer.New("my-prod-stack"),
  }
}

func (y *MyInfra) deploy() {
  name := y.NewResourceName("orders", "bucket", 64) // my-prod-stack-orders-bucket
  ...
  name = y.NewResourceName("orders", "cache", 64) // my-prod-stack-orders-bucket
  ...
  name = y.NewResourceName("pending-work", "queue", 64) // my-prod-stack-pending-work-queue
  ...
  name = y.NewResourceName("backend-processor", "service-account", 30) // my-prod-backend-pr-service-a
  ...
  name = y.NewResourceName("ingestor", "generic-service", 25) // my-prod-inges-generic-s
  ...
  name = y.NewResourceName("require-https", "", 20) // my-pro-require-https
}
```

### Install

```sh
go get github.com/davidmontoyago/commodity-namer
```


### Name structure

`"[base name]-<resource name>-[resource type or group]"`

1. **base name:**: optional prefix to set for all resources. E.g.: gcp-
2. **resource name:**: required resource name. E.g.: document-store, task-backlog, assets-cache, inference-endpoint
3. **resource type:** optional resource type or group. E.g: secret, bucket, service, version
