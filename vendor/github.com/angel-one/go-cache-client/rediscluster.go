package cache

import (
	"context"
	"crypto/tls"
	"reflect"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"

	jsoniter "github.com/json-iterator/go"
)

type redisClusterCacheClient struct {
	loader            Loader
	logger            Logger
	defaultExpiration time.Duration
	client            *redis.ClusterClient
}

type redisClusterCacheConnectionConfig struct {
	idleConnectionTimeout time.Duration
	expectContinueTimeout time.Duration
	dialTimeout           time.Duration
	readTimeout           time.Duration
	writeTimeout          time.Duration
}

type redisClusterCacheRetryConfig struct {
	maxRetries      int
	minRetryBackoff time.Duration
	maxRetryBackoff time.Duration
}

type redisClusterCachePoolConfig struct {
	poolSize           int
	minIdleConnections int
	maxConnectionAge   time.Duration
	poolTimeout        time.Duration
}

type redisClusterCacheCredentialConfig struct {
	username string
	password string
}

func getRedisClusterCacheConnectionConfig(options map[string]interface{}) redisClusterCacheConnectionConfig {
	// providing the defaults to the connection config
	config := redisClusterCacheConnectionConfig{
		idleConnectionTimeout: time.Second * 90,
		expectContinueTimeout: time.Second,
		dialTimeout:           time.Second,
		readTimeout:           time.Second,
		writeTimeout:          time.Second,
	}
	// now checking for overrides
	if val, err := getDurationConfigOption(options, "idleConnectionTimeout"); err == nil {
		config.idleConnectionTimeout = val
	}
	if val, err := getDurationConfigOption(options, "expectContinueTimeout"); err == nil {
		config.expectContinueTimeout = val
	}
	if val, err := getDurationConfigOption(options, "dialTimeout"); err == nil {
		config.dialTimeout = val
	}
	if val, err := getDurationConfigOption(options, "readTimeout"); err == nil {
		config.readTimeout = val
	}
	if val, err := getDurationConfigOption(options, "writeTimeout"); err == nil {
		config.writeTimeout = val
	}
	return config
}

func getRedisClusterCacheRetryConfig(options map[string]interface{}) redisClusterCacheRetryConfig {
	config := redisClusterCacheRetryConfig{
		maxRetries:      3,
		minRetryBackoff: time.Millisecond,
		maxRetryBackoff: 10 * time.Millisecond,
	}
	// now checking for overrides
	if val, err := getIntConfigOption(options, "maxRetries"); err == nil {
		config.maxRetries = val
	}
	if val, err := getDurationConfigOption(options, "minRetryBackoff"); err == nil {
		config.minRetryBackoff = val
	}
	if val, err := getDurationConfigOption(options, "maxRetryBackoff"); err == nil {
		config.maxRetryBackoff = val
	}
	return config
}

func getRedisClusterCachePoolConfig(options map[string]interface{}) redisClusterCachePoolConfig {
	config := redisClusterCachePoolConfig{
		poolSize:           5 * (runtime.GOMAXPROCS(0) + 1),
		minIdleConnections: 1,
		maxConnectionAge:   90 * time.Second,
		poolTimeout:        time.Second,
	}
	// now checking for overrides
	if val, err := getIntConfigOption(options, "poolSize"); err == nil {
		config.poolSize = val
	}
	if val, err := getIntConfigOption(options, "minIdleConnections"); err == nil {
		config.minIdleConnections = val
	}
	if val, err := getDurationConfigOption(options, "maxConnectionAge"); err == nil {
		config.maxConnectionAge = val
	}
	if val, err := getDurationConfigOption(options, "poolTimeout"); err == nil {
		config.poolTimeout = val
	}
	return config
}

func getRedisClusterCacheCredentialConfig(options map[string]interface{}) redisClusterCacheCredentialConfig {
	config := redisClusterCacheCredentialConfig{}
	// now checking for overrides
	if val, err := getStringConfigOption(options, "username"); err == nil {
		config.username = val
	}
	if val, err := getStringConfigOption(options, "password"); err == nil {
		config.password = val
	}
	return config
}

func getRedisClusterCacheTLSConfig(options map[string]interface{}) *tls.Config {
	var val interface{}
	var ok bool
	var c *tls.Config
	if options == nil {
		return nil
	}
	if val, ok = options["tlsConfig"]; ok {
		if c, ok = val.(*tls.Config); !ok {
			return nil
		}
	} else {
		return nil
	}
	return c
}

func getRedisClusterCacheClientOptions(options map[string]interface{}) (*redis.ClusterOptions, error) {
	c := getRedisClusterCacheConnectionConfig(options)
	r := getRedisClusterCacheRetryConfig(options)
	p := getRedisClusterCachePoolConfig(options)
	cr := getRedisClusterCacheCredentialConfig(options)
	clientOptions := &redis.ClusterOptions{
		Username:        cr.username,
		Password:        cr.password,
		MaxRetries:      r.maxRetries,
		MinRetryBackoff: r.minRetryBackoff,
		MaxRetryBackoff: r.maxRetryBackoff,
		DialTimeout:     c.dialTimeout,
		ReadTimeout:     c.readTimeout,
		WriteTimeout:    c.writeTimeout,
		PoolSize:        p.poolSize,
		MinIdleConns:    p.minIdleConnections,
		ConnMaxLifetime: p.maxConnectionAge,
		PoolTimeout:     p.poolTimeout,
		ConnMaxIdleTime: c.idleConnectionTimeout,
		TLSConfig:       getRedisClusterCacheTLSConfig(options),
	}
	// now fill the addresses
	addresses, err := getStringSliceConfigOption(options, "addresses")
	if err != nil {
		return nil, err
	}
	clientOptions.Addrs = addresses
	return clientOptions, nil
}

func getRedisClusterCacheDefaultExpiration(options map[string]interface{}) time.Duration {
	expiration := NoExpirationDuration
	if val, err := getDurationConfigOption(options, "defaultExpiration"); err == nil {
		expiration = val
	}
	return expiration
}

func (r *redisClusterCacheClient) getCMD(ctx context.Context, key string) *redis.Cmd {
	return r.client.Do(ctx, "get", key)
}

func (r *redisClusterCacheClient) getBytes(ctx context.Context, key string) ([]byte, error) {
	t, err := r.getCMD(ctx, key).Text()
	if err != nil {
		return nil, err
	}
	return []byte(t), nil
}

func (r *redisClusterCacheClient) getSliceOrMap(ctx context.Context, key string) (interface{}, error) {
	b, err := r.getBytes(ctx, key)
	if err != nil {
		return nil, err
	}
	var d interface{}
	err = jsoniter.Unmarshal(b, &d)
	return d, err
}

func newRedisClusterCacheClient(ctx context.Context, options Options) (*redisClusterCacheClient, error) {
	clientOptions, err := getRedisClusterCacheClientOptions(options.Params)
	if err != nil {
		return nil, err
	}

	r := &redisClusterCacheClient{
		loader:            options.Loader,
		logger:            options.Logger,
		client:            redis.NewClusterClient(clientOptions),
		defaultExpiration: getRedisClusterCacheDefaultExpiration(options.Params),
	}
	err = r.client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *redisClusterCacheClient) log(ctx context.Context, startTS time.Time, err error) {
	if r.logger != nil {
		r.logger(ctx, getRequiredCallerFunctionName(), time.Now().Sub(startTS).Milliseconds(), err)
	}
}

func (r *redisClusterCacheClient) Get(ctx context.Context, key string) (result interface{}, err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.getCMD(ctx, key).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = r.loader(ctx, key)
			if err != nil {
				return nil, err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.SetD(ctx, key, d)
			return d, nil
		}
		// otherwise, use the error above
		return nil, ErrNotFound
	}
	// we have the value already
	return d, err
}

func (r *redisClusterCacheClient) GetInt(ctx context.Context, key string) (int64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toInt64(d)
}

func (r *redisClusterCacheClient) GetFloat(ctx context.Context, key string) (float64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toFloat64(d)
}

func (r *redisClusterCacheClient) GetString(ctx context.Context, key string) (string, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return toString(d)
}

func (r *redisClusterCacheClient) GetBool(ctx context.Context, key string) (bool, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return toBool(d)
}

func (r *redisClusterCacheClient) GetSlice(ctx context.Context, key string) ([]interface{}, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toSlice(d)
}

func (r *redisClusterCacheClient) GetIntSlice(ctx context.Context, key string) ([]int64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Slice(d)
}

func (r *redisClusterCacheClient) GetFloatSlice(ctx context.Context, key string) ([]float64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Slice(d)
}

func (r *redisClusterCacheClient) GetStringSlice(ctx context.Context, key string) ([]string, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringSlice(d)
}

func (r *redisClusterCacheClient) GetBoolSlice(ctx context.Context, key string) ([]bool, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolSlice(d)
}

func (r *redisClusterCacheClient) GetMap(ctx context.Context, key string) (map[string]interface{}, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toMap(d)
}

func (r *redisClusterCacheClient) GetIntMap(ctx context.Context, key string) (map[string]int64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Map(d)
}

func (r *redisClusterCacheClient) GetFloatMap(ctx context.Context, key string) (map[string]float64, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Map(d)
}

func (r *redisClusterCacheClient) GetStringMap(ctx context.Context, key string) (map[string]string, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringMap(d)
}

func (r *redisClusterCacheClient) GetBoolMap(ctx context.Context, key string) (map[string]bool, error) {
	d, err := r.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolMap(d)
}

func (r *redisClusterCacheClient) GetAll(ctx context.Context, keys ...string) map[string]interface{} {
	ts := time.Now()
	defer r.log(ctx, ts, nil)
	result := make(map[string]interface{})
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)
	for _, key := range keys {
		_ = pipeline.Get(ctx, key)
	}
	// do the job
	commands, _ := pipeline.Exec(ctx)

	for i, cmd := range commands {
		value, err := cmd.(*redis.StringCmd).Result()
		if err == nil {
			result[keys[i]] = value
		}
	}
	return result
}

func (r *redisClusterCacheClient) GetJson(ctx context.Context, key string, value interface{}) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.getCMD(ctx, key).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetJsonBytes(ctx, key, r.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.SetD(ctx, key, d)
			return castFromJson(d, value)
		}
		// otherwise, use the error above
		return ErrNotFound
	}
	// we have the value already
	return castFromJson(d, value)
}

func (r *redisClusterCacheClient) GetXml(ctx context.Context, key string, value interface{}) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.getCMD(ctx, key).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetXmlBytes(ctx, key, r.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.SetD(ctx, key, d)
			return castFromXml(d, value)
		}
		// otherwise, use the error above
		return ErrNotFound
	}
	// we have the value already
	return castFromXml(d, value)
}

func (r *redisClusterCacheClient) GetYaml(ctx context.Context, key string, value interface{}) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.getCMD(ctx, key).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetYamlBytes(ctx, key, r.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.SetD(ctx, key, d)
			return castFromYaml(d, value)
		}
		// otherwise, use the error above
		return ErrNotFound
	}
	// we have the value already
	return castFromYaml(d, value)
}

func (r *redisClusterCacheClient) GetStringSet(ctx context.Context, key string) (result []string, err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.client.SMembers(ctx, key).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader exists, so in that case try to load it
			d, err = loadAndGetStringSlice(ctx, key, r.loader)
			if err != nil {
				return nil, err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.AddToStringSetD(ctx, key, d...)
			return d, nil
		}
		// otherwise, use the error above
		return nil, ErrNotFound
	}
	// we have the value already
	return d, err
}

func (r *redisClusterCacheClient) GetStringList(ctx context.Context, key string, offset, limit int64) (result []string, err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d, err := r.client.ZRange(ctx, key, offset, offset+limit).Result()
	if err == redis.Nil {
		// this means either the data is past expiry, or the data does not exist
		if r.loader != nil {
			// now that the loader exists, so in that case try to load it
			d, err = loadAndGetStringSlice(ctx, key, r.loader)
			if err != nil {
				return nil, err
			}
			// otherwise, set the data in the cache, also use it
			_ = r.AddToStringListD(ctx, key, d...)
			return d, nil
		}
		// otherwise, use the error above
		return nil, ErrNotFound
	}
	// we have the value already
	return d, err
}

func (r *redisClusterCacheClient) GetMultipleStringLists(ctx context.Context, keys ...string) (result map[string][]string, err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	// create the result
	result = make(map[string][]string)
	// begin the pipeline
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)
	// fetch the result
	for _, key := range keys {
		d, err := pipeline.ZRange(ctx, key, 0, -1).Result()
		if err == nil {
			result[key] = d
		}
	}
	// return the result
	return result, nil
}

func (r *redisClusterCacheClient) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	if expiry < 0 {
		expiry = 0
	}
	// get the kind
	kind := reflect.ValueOf(value).Kind()
	// we need to handle slices and maps particularly
	if kind == reflect.Slice || kind == reflect.Map {
		value, err = toJsonBytes(value)
		if err != nil {
			return err
		}
	}
	// otherwise, a generic handling
	return r.client.Set(ctx, key, value, expiry).Err()
}

func (r *redisClusterCacheClient) SetD(ctx context.Context, key string, value interface{}) error {
	return r.Set(ctx, key, value, r.defaultExpiration)
}

func (r *redisClusterCacheClient) SetAll(ctx context.Context, expiry time.Duration, values map[string]interface{}) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	if expiry < 0 {
		expiry = 0
	}
	// begin the pipeline
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)
	// queue the job
	for key, value := range values {
		// get the kind
		kind := reflect.ValueOf(value).Kind()
		// we need to handle slices and maps particularly
		if kind == reflect.Slice || kind == reflect.Map {
			value, err = toJsonBytes(value)
			if err != nil {
				return err
			}
		}
		// otherwise, a generic handling
		_ = pipeline.Set(ctx, key, value, expiry)
	}
	// do the job
	_, err = pipeline.Exec(ctx)
	return err
}

func (r *redisClusterCacheClient) SetAllD(ctx context.Context, values map[string]interface{}) error {
	return r.SetAll(ctx, r.defaultExpiration, values)
}

func (r *redisClusterCacheClient) SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toJsonBytes(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, b, expiry)
}

func (r *redisClusterCacheClient) SetJsonD(ctx context.Context, key string, value interface{}) error {
	return r.SetJson(ctx, key, value, r.defaultExpiration)
}

func (r *redisClusterCacheClient) SetXml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toXmlBytes(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, b, expiry)
}

func (r *redisClusterCacheClient) SetXmlD(ctx context.Context, key string, value interface{}) error {
	return r.SetXml(ctx, key, value, r.defaultExpiration)
}

func (r *redisClusterCacheClient) SetYaml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toYamlBytes(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, b, expiry)
}

func (r *redisClusterCacheClient) SetYamlD(ctx context.Context, key string, value interface{}) error {
	return r.SetYaml(ctx, key, value, r.defaultExpiration)
}

func (r *redisClusterCacheClient) AddToStringSet(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d := make([]interface{}, len(values))
	for i, v := range values {
		d[i] = v
	}
	return r.client.Watch(ctx, func(tx *redis.Tx) error {
		err := tx.SAdd(ctx, key, d...).Err()
		if err != nil {
			return err
		}
		if expiry > 0 {
			err = tx.Expire(ctx, key, expiry).Err()
			if err != nil {
				return err
			}
		}
		return nil
	}, key)
}

func (r *redisClusterCacheClient) AddToStringSetD(ctx context.Context, key string, values ...string) error {
	return r.AddToStringSet(ctx, key, r.defaultExpiration, values...)
}

func (r *redisClusterCacheClient) AddToStringList(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d := make([]redis.Z, len(values))
	for idx, v := range values {
		d[idx] = redis.Z{
			Score:  1, // all elements will have same score
			Member: v,
		}
	}
	return r.client.Watch(ctx, func(tx *redis.Tx) error {
		err := tx.ZAdd(ctx, key, d...).Err()
		if err != nil {
			return err
		}
		if expiry > 0 {
			err = tx.Expire(ctx, key, expiry).Err()
			if err != nil {
				return err
			}
		}
		return nil
	}, key)
}

func (r *redisClusterCacheClient) AddToStringListD(ctx context.Context, key string, values ...string) error {
	return r.AddToStringList(ctx, key, r.defaultExpiration, values...)
}

func (r *redisClusterCacheClient) IncrementAll(ctx context.Context, expiry time.Duration, keys ...string) (result map[string]int64, err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	// create the result
	result = make(map[string]int64)
	// queue the commands
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	var incrResult []*redis.IntCmd
	// fetch the result
	for _, key := range keys {
		incrResult = append(incrResult, pipeline.Incr(ctx, key))
		if err != nil {
			return nil, err
		}

		if expiry > 0 {
			// Set expiration for each key
			if err := pipeline.Expire(ctx, key, expiry).Err(); err != nil {
				return result, err
			}
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	// fetch the results
	for i, cmd := range incrResult {
		d, err := cmd.Result()
		if err == nil {
			result[keys[i]] = d
		}
	}
	// return the result
	return result, nil
}

func (r *redisClusterCacheClient) IncrementIntAllD(ctx context.Context, keys ...string) (map[string]int64, error) {
	return r.IncrementAll(ctx, r.defaultExpiration, keys...)
}

func (r *redisClusterCacheClient) IncrementWithExpiry(ctx context.Context, key string,
	expiry time.Duration) (result int64, err error) {

	ts := time.Now()
	defer r.log(ctx, ts, err)
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	//to increment value
	incrResult := pipeline.Incr(ctx, key)
	if err != nil {
		return result, err
	}

	//update expiry
	if err := pipeline.Expire(ctx, key, expiry).Err(); err != nil {
		return result, err
	}

	//execute the commands
	if _, err := pipeline.Exec(ctx); err != nil {
		return result, err
	}

	//fetch the new value post increment
	result, err = incrResult.Result()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *redisClusterCacheClient) IncrementWithGreaterExpiry(ctx context.Context, key string, value int64,
	expiry time.Duration) (result int64, err error) {

	ts := time.Now()
	defer r.log(ctx, ts, err)
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	//to increment value
	incrResult := pipeline.IncrBy(ctx, key, value)
	if err != nil {
		return result, err
	}

	if expiry != 0 {
		// add expiry if expire is not set to the key
		if err := pipeline.ExpireNX(ctx, key, expiry).Err(); err != nil {
			return result, err
		}
	}

	// update expiry if current expire is greater than old ttl of the key
	if err := pipeline.ExpireGT(ctx, key, expiry).Err(); err != nil {
		return result, err
	}

	//execute the commands
	if _, err := pipeline.Exec(ctx); err != nil {
		return result, err
	}

	//fetch the new value post increment
	result, err = incrResult.Result()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *redisClusterCacheClient) SetExpiryAll(ctx context.Context, expiry time.Duration, keys ...string) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	// queue the commands
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)
	for _, key := range keys {
		_, err := pipeline.Expire(ctx, key, expiry).Result()
		if err != nil {
			return err
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClusterCacheClient) Delete(ctx context.Context, keys ...string) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	return r.client.Del(ctx, keys...).Err()
}
func (r *redisClusterCacheClient) DeleteAll(ctx context.Context) (err error) {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisClusterCacheClient) DeleteFromStringSet(ctx context.Context, key string, values ...string) (err error) {
	ts := time.Now()
	defer r.log(ctx, ts, err)
	d := make([]interface{}, len(values))
	for i, v := range values {
		d[i] = v
	}
	c, err := r.client.SRem(ctx, key, d...).Result()
	if err != nil {
		return ErrIncorrectValueKind
	}
	if c == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *redisClusterCacheClient) DeleteFromStringList(ctx context.Context, key string, values ...string) (err error) {
	d := make([]interface{}, len(values))
	for i, v := range values {
		d[i] = v
	}
	c, err := r.client.ZRem(ctx, key, d...).Result()
	if err != nil {
		return ErrIncorrectValueKind
	}
	if c == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *redisClusterCacheClient) Subscribe(ctx context.Context, subscriber Subscriber, options SubscriberOptions,
	patterns ...string) SubscriberCloser {
	if subscriber == nil {
		return nil
	}
	options.init()
	ps := r.client.PSubscribe(ctx, patterns...)
	go func() {
		ch := ps.Channel(redis.WithChannelSize(options.MaxConcurrency), redis.WithChannelSendTimeout(options.ProcessTimeout))
		for msg := range ch {
			subscriber(ctx, msg.Channel, msg.Payload)
		}
	}()
	return ps.Close
}

func (r *redisClusterCacheClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *redisClusterCacheClient) Close() error {
	return r.client.Close()
}

func (r *redisClusterCacheClient) closePipeline(p redis.Pipeliner) {}

func (r *redisClusterCacheClient) GetMultipleHashStringMap(ctx context.Context, keys ...string) (map[string]map[string]string, error) {
	ts := time.Now()
	var err error
	defer r.log(ctx, ts, err)
	// create the result
	result := make(map[string]map[string]string)
	// queue the commands
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)
	// fetch the result
	for _, key := range keys {
		data := pipeline.HGetAll(ctx, key)
		err = data.Err()
		if err != nil {
			return nil, err
		}
	}
	commands, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	// fetch the results
	for i, cmd := range commands {
		d, err := cmd.(*redis.MapStringStringCmd).Result()
		if err == nil {
			result[keys[i]] = d
		}
	}
	// return the result
	return result, nil
}

func (r *redisClusterCacheClient) SetHashStringMap(ctx context.Context, key string, value []string) error {
	return r.client.HSet(ctx, key, value).Err()
}

func (r *redisClusterCacheClient) DeleteFromHashStringMap(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *redisClusterCacheClient) SetHashStringMapWithGreaterExpiry(ctx context.Context, key string, value []string, expiry time.Duration) error {

	ts := time.Now()
	var err error
	defer r.log(ctx, ts, err)
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	//to set value
	hSetResult := pipeline.HSet(ctx, key, value)
	if hSetResult.Err() != nil {
		return hSetResult.Err()
	}

	if expiry != 0 {
		// add expiry if expire is not set to the key
		if err := pipeline.ExpireNX(ctx, key, expiry).Err(); err != nil {
			return err
		}
	}

	// update expiry if current expire is greater than old ttl of the key
	if err := pipeline.ExpireGT(ctx, key, expiry).Err(); err != nil {
		return err
	}

	//execute the commands
	if _, err := pipeline.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *redisClusterCacheClient) SetHashStringMapWithGreaterExpiryAndDelete(ctx context.Context, key, deleteKey string, value, deleteFields []string, expiry time.Duration) error {

	ts := time.Now()
	var err error
	defer r.log(ctx, ts, err)
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	//to set value
	hSetResult := pipeline.HSet(ctx, key, value)
	if hSetResult.Err() != nil {
		return hSetResult.Err()
	}

	if expiry != 0 {
		// add expiry if expire is not set to the key
		if err := pipeline.ExpireNX(ctx, key, expiry).Err(); err != nil {
			return err
		}
	}

	// update expiry if current expire is greater than old ttl of the key
	if err := pipeline.ExpireGT(ctx, key, expiry).Err(); err != nil {
		return err
	}

	//delete
	if err := pipeline.HDel(ctx, deleteKey, deleteFields...).Err(); err != nil {
		return err
	}

	//execute the commands
	if _, err := pipeline.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *redisClusterCacheClient) DeleteFromMultipleHashStringMap(ctx context.Context, keyValueHashMap map[string][]string) error {
	ts := time.Now()
	var err error
	defer r.log(ctx, ts, err)
	pipeline := r.client.Pipeline()
	defer r.closePipeline(pipeline)

	for key, val := range keyValueHashMap {
		if err := pipeline.HDel(ctx, key, val...).Err(); err != nil {
			return err
		}
	}
	//execute the commands
	if _, err := pipeline.Exec(ctx); err != nil {
		return err
	}

	return nil
}
