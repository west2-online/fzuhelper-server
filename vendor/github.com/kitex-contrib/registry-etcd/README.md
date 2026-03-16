# registry-etcd (*This is a community driven project*)

Some application runtime use [etcd](https://github.com/etcd-io/etcd) for service discovery.

## How to use with Kitex server?

```go
import (
    ...
    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
    ...
)

func main() {
    ...
    r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"}) // r should not be reused.
    if err != nil {
        log.Fatal(err)
    }
    // https://www.cloudwego.io/docs/tutorials/framework-exten/registry/#integrate-into-kitex
    server, err := echo.NewServer(new(EchoImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}), server.WithRegistry(r))
    if err != nil {
        log.Fatal(err)
    }
    err = server.Run()
    if err != nil {
        log.Fatal(err)
    }
    ...
}
```


## How to use with Kitex client?

```go
import (
    ...
    "github.com/cloudwego/kitex/client"
    etcd "github.com/kitex-contrib/registry-etcd"
    ...
)

func main() {
    ...
    r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
    if err != nil {
        log.Fatal(err)
    }
    client, err := echo.NewClient("echo", client.WithResolver(r))
    if err != nil {
        log.Fatal(err)
    }
    ...
}
```

## Authentication

### server
```go
package main

import (
    ...
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

type HelloImpl struct{}

func (h *HelloImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	resp = &api.Response{
		Message: req.Message,
	}
	return
}

func main() {
	// creates a etcd based registry with given username and password
	r, err := etcd.NewEtcdRegistryWithAuth([]string{"127.0.0.1:2379"}, "username", "password")
	if err != nil {
		log.Fatal(err)
	}
	server := hello.NewServer(new(HelloImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "Hello",
	}))
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

```

### client
```go
package main

import (
    ...
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() { 
	// creates a etcd based resolver with given username and password
	r, err := etcd.NewEtcdResolverWithAuth([]string{"127.0.0.1:2379"}, "username", "password")
	if err != nil {
		log.Fatal(err)
	}
	client := hello.MustNewClient("Hello", client.WithResolver(r))
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		resp, err := client.Echo(ctx, &api.Request{Message: "Hello"})
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}
```

## Retry

After the service is registered to ETCD, it will regularly check the status of the service. If any abnormal status is found, it will try to register the service again. `ObserveDelay` is the delay time for checking the service status under normal conditions, and `RetryDelay` is the delay time for attempting to register the service after disconnecting.

### Default Retry Config

| Config Name         | Default Value    | Description                                                                               |
|:--------------------|:-----------------|:------------------------------------------------------------------------------------------|
| WithMaxAttemptTimes | 5                | Used to set the maximum number of attempts, if 0, it means infinite attempts              |
| WithObserveDelay    | 30 * time.Second | Used to set the delay time for checking service status under normal connection conditions |
| WithRetryDelay      | 10 * time.Second | Used to set the retry delay time after disconnecting                                      |

### Example

If you do not need to customize the retry configuration, use `etcd. NewEtcdRegistry()`.

If you need to customize the retry configuration, use the following code:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api"
	"github.com/cloudwego/kitex-examples/hello/kitex_gen/api/hello"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
)

type HelloImpl struct{}

func (h *HelloImpl) Echo(ctx context.Context, req *api.Request) (resp *api.Response, err error) {
	resp = &api.Response{
		Message: req.Message,
	}
	return
}

func main() {
	retryConfig := retry.NewRetryConfig(
		retry.WithMaxAttemptTimes(10),
		retry.WithObserveDelay(20*time.Second),
		retry.WithRetryDelay(5*time.Second),
	)
	r, err := etcd.NewEtcdRegistryWithRetry([]string{"127.0.0.1:2379"}, retryConfig)
	if err != nil {
		log.Fatal(err)
	}
	server := hello.NewServer(new(HelloImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "Hello",
	}))
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Default Weight

The weighted load balancing algorithm can only handle positive weights, and will be filtered when the weight is if 0 or negative. Setting default weights can avoid filtering.

### Default Config

| Config Name       | Default Value | Description                                                                                                                          |
|:------------------|:--------------|:-------------------------------------------------------------------------------------------------------------------------------------|
| WithDefaultWeight | 10            | Used to set the default wight of instances, if 0 or negative, it means instances with 0 or negative weight will be filtered          |

### Example

```go
package main

import (
	...
    "github.com/cloudwego/kitex/client"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	// creates a etcd based resolver with default weight
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"}, etcd.WithDefaultWeight(10))
	if err != nil {
		log.Fatal(err)
	}
	client := hello.MustNewClient("Hello", client.WithResolver(r))
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		resp, err := client.Echo(ctx, &api.Request{Message: "Hello"})
		cancel()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp)
		time.Sleep(time.Second)
	}
}
```

## How to Dynamically specify ip and port
To dynamically specify an IP and port, one should first set the environment variables KITEX_IP_TO_REGISTRY and KITEX_PORT_TO_REGISTRY. If these variables are not set, the system defaults to using the service's listening IP and port. Notably, if the service's listening IP is either not set or set to "::", the system will automatically retrieve and use the machine's IPV4 address.

## More info

See [example](/example).


## Compatibility

Compatible with etcd server v3.




maintained by: [lizheming](https://github.com/duduainankai)
