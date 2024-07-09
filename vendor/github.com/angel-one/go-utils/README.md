# Go Utils

[![Go Analyser](http://ec2-3-6-93-149.ap-south-1.compute.amazonaws.com:8080/badge/github.com/angel-one/go-utils)](http://ec2-3-6-93-149.ap-south-1.compute.amazonaws.com:8080/report/github.com/angel-one/go-utils)
[![Test Coverage](http://ec2-3-6-93-149.ap-south-1.compute.amazonaws.com:8080/testCoverageBadge/github.com/angel-one/go-utils)](http://ec2-3-6-93-149.ap-south-1.compute.amazonaws.com:8080/report/github.com/angel-one/go-utils)

Go utils contains the set of common reusable utility methods which can be reused across all Angel One Go projects.

## Project Versioning

Go utils uses [semantic versioning](http://semver.org/). API should not change between patch and minor releases. New minor versions may add additional features to the API.

## Installation

To install `Go Utils` package, you need to install Go and set your Go workspace first.

1. The first need Go installed (version 1.12+ is required), then you can use the below Go command to install Go Config Client.

```shell
go get github.com/angel-one/go-utils
```

2. Because this is a private repository, you will need to mark this in the `Go env variables`.

```shell
go env -w GOPRIVATE=github.com/angel-one/go-utils
```

3. Also, follow [this](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) to generate a personal access token and add the following line to your `$HOME/.netrc` file.

```
machine github.com login ${USERNAME} password ${PERSONAL_ACCESS_TOKEN}
```

4. Import it in your code:

```go
package sample

import "github.com/angel-one/go-utils"
```

## Usage

### Utils

```go
package sample

import (
	"github.com/angel-one/go-utils"
	"io/ioutil"
	"strings"
)

type Sample struct {
	A string `json:"a"`
}

func sample() {
	r := ioutil.NopCloser(strings.NewReader("{\"a\": \"naruto\"}"))
	
	// get data as bytes
	b, err := utils.GetDataAsBytes(r)

	// get data as string
	s, err := utils.GetDataAsString(r)
	
	// get data unmarshalled to struct
	v := Sample{}
	err = utils.GetJSONData(r, &v)
}
```

### Logger

The logger methods need to share the context to log the unique identifier for the context, to be able to be clubbed together.

```go
package sample

import (
	"context"
	utilConstants "github.com/angel-one/go-utils/constants"
	"github.com/angel-one/go-utils/log"
)

func init() {
	// used to initialize logger
	log.InitLogger(utilConstants.TraceLevel)
}

func sample() {
	ctx := context.Background()
	log.Trace(ctx).Msg("sample trace log")
	log.Debug(ctx).Msg("sample debug log")
	log.Info(ctx).Msg("sample info log")
	log.Warn(ctx).Msg("sample warn log")
	log.Error(ctx).Msg("sample error log")
	log.Fatal(ctx).Msg("sample fatal log")
	log.Panic(ctx).Msg("sample panic log")
	
	// if there is no context needed to be used
	// then even nil can be passed safely
	log.Trace(nil).Msg("sample trace log")
	log.Debug(nil).Msg("sample debug log")
	log.Info(nil).Msg("sample info log")
	log.Warn(nil).Msg("sample warn log")
	log.Error(nil).Msg("sample error log")
	log.Fatal(nil).Msg("sample fatal log")
	log.Panic(nil).Msg("sample panic log")
}
```

### Logger Middleware

Logger GIN middleware is used to log the request details and add a unique identifier to the context which will be sent with every log associated with the request context.

The identifier can either be sent as a request header, which has to be with the key `X-requestId`, otherwise it is generated uniquely.

```go
package sample

import (
	"github.com/angel-one/go-utils/middlewares"
	"github.com/gin-gonic/gin"
)

func Sample() {
	router := gin.New()
	router.Use(middlewares.Logger(middlewares.LoggerMiddlewareOptions{}))
	
	_ = router.Run(":8080")
}
```

### Supported logging context parameters
The logging library automatically extracts the following fields from the context and logs them if set:

- `id`: any ID for the log
- `path`: API Path
- `correlationId`: correlation ID to correlate multiple logs for the same flow
- `clientID`: client ID of the user
- `action`: action for which the log was captured

These are present in `constants/logparams.go`` as constants for reuse.