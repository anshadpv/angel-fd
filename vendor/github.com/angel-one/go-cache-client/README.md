# Go Cache Client

This is the config client for Go projects.

- **In-Memory Cache Client** - Fully Implemented
- **Redis Cache Client** - Fully Implemented
- **Redis Cluster Cache Client** - Fully Implemented
- **Memcache Client** - TBD

## Project Versioning

Go Cache Client uses [semantic versioning](http://semver.org/). API should not change between patch and minor releases. New minor versions may add additional features to the API.

## Installation

To install `Go Cache Client` package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.13+ is required), then you can use the below Go command to install Go Config Client.

```shell
go get github.com/angel-one/go-cache-client
```

2. Because this is a private repository, you will need to mark this in the Go env variables.

```shell
go env -w GOPRIVATE=github.com/angel-one/go-cache-client
```

3. Also, follow this to generate a personal access token and add the following line to your $HOME/.netrc file.

```
machine github.com login ${USERNAME} password ${PERSONAL_ACCESS_TOKEN}
```

4. Import it in your code:

```go
import cache "github.com/angel-one/go-cache-client"
```

## Usage

### New Client

```go
import (
	"context"
	cache "github.com/angel-one/go-cache-client"
	"time"
)

inMemoryCacheClient, err := configs.New(context.Background(), configs.Options{
    Provider: cache.InMemory,
    Params: map[string]interface{}{
        "defaultExpiration": time.Minute * 5,
    },
})
if err != nil {
	// handle error
}
```

### Getting Data from Cache

There are several typed methods available which can be used to get the data from cache.

### Setting Data in Cache

There are 2 types of methods available.
1. Plain method which can be used to set the data with default expiration provided for the client.
2. Method which takes an extra parameter called the expiry duration.

### Deleting Data in Cache

It is possible to delete multiple keys at a time.
