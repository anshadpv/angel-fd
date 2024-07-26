package cache

import (
	"context"
	"time"

	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	cache "github.com/angel-one/go-cache-client"
)

var redisClient cache.Client

func Init(ctx context.Context, cfg *config.Client, key string) error {
	var err error
	if err = initRedis(ctx, cfg, key, constants.RedisConfigKey, nil); err != nil {
		return err
	}
	return nil
}

func initRedis(ctx context.Context, cfg *config.Client, key string, redisConfigKey string, loader func(ctx context.Context, key string) (interface{}, error)) error {
	var err error
	redisConfig, err := cfg.GetMap(key, redisConfigKey)
	if err != nil {
		log.Fatal(ctx).Err(err).Msgf("error during load of key: %s; redisConfigKey: %s configs", key, redisConfigKey)
		return err
	}
	if len(redisConfig) == 0 {
		log.Fatal(ctx).Err(err).Msgf("missing config for key: %s; redisConfigKey: %s configs", key, redisConfigKey)
		return err
	}

	redisConnections, err := cfg.GetMap(key, constants.RedisConnections)
	if err != nil {
		log.Fatal(ctx).Err(err).Msgf("unable to initialize redis connections")
		return err
	}

	redisClient, err = cache.New(ctx, cache.Options{
		Provider: cache.RedisCluster,
		Loader:   loader,
		Params: map[string]interface{}{
			constants.DefaultExpiration:                 time.Hour,
			constants.RedisAddressesKey:                 []string{redisConfig[constants.RedisUrlKey].(string)},
			constants.RedisPoolSizeKey:                  redisConnections[constants.RedisPoolSize],
			constants.RedisMaxConnAgeKey:                redisConnections[constants.RedisMaxConnAge],
			constants.RedisMinIdleConnKey:               redisConnections[constants.RedisMinIdleConn],
			constants.RedisPoolTimeoutInMilliSecondKey:  redisConnections[constants.RedisPoolTimeoutInMilliSecond],
			constants.RedisReadTimeoutInMilliSecondKey:  redisConnections[constants.RedisReadTimeoutInMilliSecond],
			constants.RedisWriteTimeoutInMilliSecondKey: redisConnections[constants.RedisWriteTimeoutInMilliSecond],
		},
	})
	if err != nil {
		log.Fatal(ctx).Err(err).Msgf("failed to initlaize redis cache")
		return err
	}
	log.Info(ctx).Msg("initialised redis cache client")
	return nil
}
func GetRedisClient() cache.Client {
	return redisClient
}
