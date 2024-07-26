package cache

import (
	"context"
	"time"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
)

// Errors
const (
	funNotImplementedError = "function not implemented yet, please implement and use"
)

type lruInMemoryCacheClientOptions struct {
	defaultExpiration time.Duration
	size              int
}

func getLRUInMemoryCacheClientOptions(options map[string]interface{}) lruInMemoryCacheClientOptions {
	var clientOptions lruInMemoryCacheClientOptions
	var err error
	clientOptions.defaultExpiration, err = getDurationConfigOption(options, "defaultExpiration")
	if err != nil {
		clientOptions.defaultExpiration = NoExpirationDuration
	}
	clientOptions.size, err = getIntConfigOption(options, "capacity")
	if err != nil {
		clientOptions.size = 2147483647
	}
	return clientOptions
}

func newLRUInMemoryCacheClient(_ context.Context, options Options) *lruInMemoryCacheClient {
	clientOptions := getLRUInMemoryCacheClientOptions(options.Params)
	i := &lruInMemoryCacheClient{
		signal:  make(chan struct{}, 1),
		cache:   lru.NewLRU[string, interface{}](clientOptions.size, nil, clientOptions.defaultExpiration),
		loader:  options.Loader,
		logger:  options.Logger,
		options: clientOptions,
	}

	return i
}

type lruInMemoryCacheClient struct {
	signal  chan struct{}
	cache   *lru.LRU[string, interface{}]
	loader  Loader
	logger  Logger
	options lruInMemoryCacheClientOptions
}

func (l *lruInMemoryCacheClient) Get(ctx context.Context, key string) (result interface{}, err error) {
	result, ok := l.cache.Get(key)
	if !ok {
		// this means either the data is past expiry, or the data does not exist
		if l.loader != nil {
			// now get data from provider
			newData, err := l.loader(ctx, key)
			if err != nil {
				return nil, err
			}

			//set the data in the cache, also use it
			_ = l.SetD(ctx, key, newData)
			return newData, nil
		}
		// otherwise, use the error above
		return nil, err
	}
	// we have the value already
	return result, nil
}

func (l *lruInMemoryCacheClient) Set(_ context.Context, key string, value interface{}, _ time.Duration) (err error) {
	l.cache.Add(key, value)
	return nil
}

func (l *lruInMemoryCacheClient) SetD(ctx context.Context, key string, value interface{}) (err error) {
	ts := time.Now()
	defer l.log(ctx, ts, err)
	l.Set(ctx, key, value, -1)
	return nil
}

func (l *lruInMemoryCacheClient) GetInt(ctx context.Context, key string) (int64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toInt64(d)
}

func (l *lruInMemoryCacheClient) GetFloat(ctx context.Context, key string) (float64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toFloat64(d)
}

func (l *lruInMemoryCacheClient) GetString(ctx context.Context, key string) (string, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return toString(d)
}

func (l *lruInMemoryCacheClient) GetBool(ctx context.Context, key string) (bool, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return toBool(d)
}

func (l *lruInMemoryCacheClient) GetSlice(ctx context.Context, key string) ([]interface{}, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toSlice(d)
}

func (l *lruInMemoryCacheClient) GetIntSlice(ctx context.Context, key string) ([]int64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Slice(d)
}

func (l *lruInMemoryCacheClient) GetFloatSlice(ctx context.Context, key string) ([]float64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Slice(d)
}

func (l *lruInMemoryCacheClient) GetStringSlice(ctx context.Context, key string) ([]string, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringSlice(d)
}

func (l *lruInMemoryCacheClient) GetBoolSlice(ctx context.Context, key string) ([]bool, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolSlice(d)
}

func (l *lruInMemoryCacheClient) GetMap(ctx context.Context, key string) (map[string]interface{}, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toMap(d)
}

func (l *lruInMemoryCacheClient) GetIntMap(ctx context.Context, key string) (map[string]int64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Map(d)
}

func (l *lruInMemoryCacheClient) GetFloatMap(ctx context.Context, key string) (map[string]float64, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Map(d)
}

func (l *lruInMemoryCacheClient) GetStringMap(ctx context.Context, key string) (map[string]string, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringMap(d)
}

func (l *lruInMemoryCacheClient) GetBoolMap(ctx context.Context, key string) (map[string]bool, error) {
	d, err := l.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolMap(d)
}

func (l *lruInMemoryCacheClient) GetAll(ctx context.Context, keys ...string) map[string]interface{} {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) GetJson(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) GetXml(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) GetYaml(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) GetStringSet(ctx context.Context, key string) ([]string, error) {
	panic(funNotImplementedError)
	return nil, nil
}

func (l *lruInMemoryCacheClient) GetStringList(ctx context.Context, key string, offset, limit int64) ([]string, error) {
	panic(funNotImplementedError)
	return nil, nil
}

func (l *lruInMemoryCacheClient) GetMultipleStringLists(ctx context.Context, keys ...string) (map[string][]string, error) {
	panic(funNotImplementedError)
	return nil, nil
}

func (l *lruInMemoryCacheClient) SetAll(ctx context.Context, expiry time.Duration, values map[string]interface{}) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetAllD(ctx context.Context, values map[string]interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetJsonD(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetXml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetXmlD(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetYaml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) SetYamlD(ctx context.Context, key string, value interface{}) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) AddToStringSet(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) AddToStringSetD(ctx context.Context, key string, values ...string) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) AddToStringList(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) AddToStringListD(ctx context.Context, key string, values ...string) error {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) IncrementAll(ctx context.Context, expiry time.Duration, keys ...string) (result map[string]int64, err error) {
	panic(funNotImplementedError)
	return nil, nil
}

func (l *lruInMemoryCacheClient) IncrementIntAllD(ctx context.Context, keys ...string) (map[string]int64, error) {
	panic(funNotImplementedError)
	return nil, nil
}

func (l *lruInMemoryCacheClient) IncrementWithExpiry(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return 0, nil
}

func (l *lruInMemoryCacheClient) IncrementWithGreaterExpiry(_ context.Context, _ string, _ int64, _ time.Duration) (int64, error) {
	return 0, nil
}

func (l *lruInMemoryCacheClient) SetExpiryAll(ctx context.Context, expiry time.Duration, keys ...string) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) Delete(ctx context.Context, keys ...string) (err error) {
	ts := time.Now()
	defer l.log(ctx, ts, err)
	for _, k := range keys {
		l.cache.Remove(k)
	}
	return nil
}

func (l *lruInMemoryCacheClient) DeleteAll(ctx context.Context) (err error) {
	ts := time.Now()
	defer l.log(ctx, ts, err)
	l.cache.Purge()
	return nil
}

func (l *lruInMemoryCacheClient) DeleteFromStringSet(ctx context.Context, key string, values ...string) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) DeleteFromStringList(ctx context.Context, key string, values ...string) (err error) {
	panic(funNotImplementedError)
	return nil
}

func (l *lruInMemoryCacheClient) Subscribe(_ context.Context, _ Subscriber, _ SubscriberOptions, _ ...string) SubscriberCloser {
	return nil
}

func (l *lruInMemoryCacheClient) Ping(ctx context.Context) error {
	return nil
}

func (l *lruInMemoryCacheClient) Close() error {
	l.signal <- struct{}{}
	l.cache.Purge()
	return nil
}

func (l *lruInMemoryCacheClient) log(ctx context.Context, startTS time.Time, err error) {
	if l.logger != nil {
		l.logger(ctx, getRequiredCallerFunctionName(), time.Now().Sub(startTS).Milliseconds(), err)
	}
}

func (l *lruInMemoryCacheClient) GetMultipleHashStringMap(ctx context.Context, keys ...string) (map[string]map[string]string, error) {
	return nil, nil
}

func (l *lruInMemoryCacheClient) SetHashStringMap(ctx context.Context, key string, value []string) error {
	return nil
}

func (l *lruInMemoryCacheClient) DeleteFromHashStringMap(ctx context.Context, key string, fields ...string) error {
	return nil
}

func (l *lruInMemoryCacheClient) SetHashStringMapWithGreaterExpiry(ctx context.Context, key string, value []string, expiry time.Duration) error {
	return nil
}

func (l *lruInMemoryCacheClient) SetHashStringMapWithGreaterExpiryAndDelete(ctx context.Context, key, deleteKey string, value, deleteFields []string, expiry time.Duration) error {
	return nil
}

func (l *lruInMemoryCacheClient) DeleteFromMultipleHashStringMap(ctx context.Context, keyValueHashMap map[string][]string) error {
	return nil
}
