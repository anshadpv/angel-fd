package httpclient

import (
	"context"
	"net/http"

	"github.com/angel-one/fd-core/commons/config"
	fderr "github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	httpclient "github.com/angel-one/go-http-client"
)

var configNotFound = fderr.New().Code("HTTP-CLIENT-01").Msg("configs not found from config file").Build()

type Client interface {
	Request(request *httpclient.Request) (*http.Response, error)
}

var client *httpclient.Client
var requestConfigs []*httpclient.RequestConfig

func NewRequest(name string) *httpclient.Request {
	return httpclient.NewRequest(name)
}

func Init(ctx context.Context, cfg *config.Client, key string, clientConfigKeys []string) error {
	log.Info(ctx).Msg("initializing http client")

	var err error
	for _, v := range clientConfigKeys {
		var config map[string]interface{}
		log.Info(ctx).Msgf("loading configs for key: %s", v)
		if config, err = cfg.GetMap(key, v); err != nil {
			log.Fatal(ctx).Err(err).Msgf("error during load of %s configs", v)
			return err
		}
		if len(config) == 0 {
			return configNotFound
		}
		requestConfig := httpclient.NewRequestConfig(v, config)
		requestConfigs = append(requestConfigs, requestConfig)
	}

	// collate all http client configs to configurer
	client = httpclient.ConfigureHTTPClient(requestConfigs...)

	return nil
}

func Default() Client {
	return client
}
