# commodity-namer

Consistent structured naming for infrastructure resources.

E.g.:

```
import "github.com/davidmontoyago/commodity-namer"

type MyInfra struct {
  // Embed the namer
	Namer
}

func NewMyInfra() *MyInfra {
  return &MyInfra{
    // Set the base name
    Namer: *namer.New("my-prod-stack"),
  }
}

func (y *MyInfra) deploy() {
  name := y.NewResourceName("orders", "bucket") // my-prod-stack-orders-bucket
  ...
  name = y.NewResourceName("orders", "cache") // my-prod-stack-orders-bucket
  ...
  name = y.NewResourceName("pending-work", "queue") // my-prod-stack-pending-work-queue
  ...
  name = y.NewResourceName("backend-processor", "service-account") // TODO
}
```

### Getting started

Embed as a type:

```go

```

Get a new name:
```go

```

### Name structure

"[base name]-<resource name>-[resource type or group]"

1. **base name:**: optional prefix to set for all resources. E.g.: gcp-
2. **resource name:**: required resource name. E.g.:
3. **resource type:** optional resource type or group.

E.g.:

```go

```
