package constants

const (
	RedisConfigKey   = "redis-config"
	RedisConnections = "redis-config.connections"

	RedisUrlKey                       = "url"
	RedisAddressesKey                 = "addresses"
	RedisPoolSizeKey                  = "poolSize"
	RedisMaxConnAgeKey                = "maxConnectionAge"
	RedisMinIdleConnKey               = "minIdleConnections"
	RedisPoolTimeoutInMilliSecondKey  = "poolTimeout"
	RedisReadTimeoutInMilliSecondKey  = "readTimeout"
	RedisWriteTimeoutInMilliSecondKey = "writeTimeout"

	DefaultExpiration              = "defaultExpiration"
	RedisPoolSize                  = "poolSize"
	RedisMaxConnAge                = "maxConn"
	RedisMinIdleConn               = "minIdleConn"
	RedisPoolTimeoutInMilliSecond  = "poolTimeoutInMilliSecond"
	RedisReadTimeoutInMilliSecond  = "readTimeoutInMilliSecond"
	RedisWriteTimeoutInMilliSecond = "redisWriteTimeoutInMilliSecond"
)
