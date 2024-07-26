package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// various modes of cache client
const (
	InMemory = iota
	Redis
	RedisCluster
	LRUInMemory
)

// Loader is the function for loading cache values on expiry automatically
// This helps avoid unnecessary conditional code that one has to write
type Loader func(ctx context.Context, key string) (interface{}, error)

// Logger is used to log the latency and error for a cache method
type Logger func(ctx context.Context, name string, latencyInMillis int64, err error)

// Subscriber is used to subscribe to events
type Subscriber func(ctx context.Context, key string, data string)

// SubscriberCloser is used to close the subscriber
type SubscriberCloser func() error

// Options is the set of configurable parameters required for initialisation of the config client
type Options struct {
	// This is the provider to be used
	Provider int
	// This is the loader for loading cache on expiry
	Loader Loader
	// These are the set of parameters required for the initialisation of the chosen parameter
	Params map[string]interface{}
	// This is the logger to be used to log various details for the cache
	Logger Logger
}

// SubscriberOptions is the set of parameters for subscription
type SubscriberOptions struct {
	// This is the maximum number of events that can be received concurrently
	// Default value is 1.
	MaxConcurrency int
	// This is the time within which the event should be processed.
	// Default value is 5 seconds.
	ProcessTimeout time.Duration
}

// Client is the contract that can be used and will be followed by every implementation of the cache client
type Client interface {
	Get(ctx context.Context, key string) (interface{}, error)
	GetInt(ctx context.Context, key string) (int64, error)
	GetFloat(ctx context.Context, key string) (float64, error)
	GetString(ctx context.Context, key string) (string, error)
	GetBool(ctx context.Context, key string) (bool, error)
	GetSlice(ctx context.Context, key string) ([]interface{}, error)
	GetIntSlice(ctx context.Context, key string) ([]int64, error)
	GetFloatSlice(ctx context.Context, key string) ([]float64, error)
	GetStringSlice(ctx context.Context, key string) ([]string, error)
	GetBoolSlice(ctx context.Context, key string) ([]bool, error)
	GetMap(ctx context.Context, key string) (map[string]interface{}, error)
	GetIntMap(ctx context.Context, key string) (map[string]int64, error)
	GetFloatMap(ctx context.Context, key string) (map[string]float64, error)
	GetStringMap(ctx context.Context, key string) (map[string]string, error)
	GetBoolMap(ctx context.Context, key string) (map[string]bool, error)

	GetAll(ctx context.Context, keys ...string) map[string]interface{}

	GetJson(ctx context.Context, key string, value interface{}) error
	GetXml(ctx context.Context, key string, value interface{}) error
	GetYaml(ctx context.Context, key string, value interface{}) error

	GetStringSet(ctx context.Context, key string) ([]string, error)
	GetStringList(ctx context.Context, key string, offset, limit int64) ([]string, error)

	GetMultipleStringLists(ctx context.Context, keys ...string) (map[string][]string, error)

	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	SetD(ctx context.Context, key string, value interface{}) error

	SetAll(ctx context.Context, expiry time.Duration, values map[string]interface{}) error
	SetAllD(ctx context.Context, values map[string]interface{}) error

	SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	SetJsonD(ctx context.Context, key string, value interface{}) error
	SetXml(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	SetXmlD(ctx context.Context, key string, value interface{}) error
	SetYaml(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	SetYamlD(ctx context.Context, key string, value interface{}) error

	AddToStringSet(ctx context.Context, key string, expiry time.Duration, values ...string) error
	AddToStringSetD(ctx context.Context, key string, values ...string) error

	AddToStringList(ctx context.Context, key string, expiry time.Duration, values ...string) error
	AddToStringListD(ctx context.Context, key string, values ...string) error

	IncrementAll(ctx context.Context, expiry time.Duration, keys ...string) (map[string]int64, error)
	IncrementIntAllD(ctx context.Context, keys ...string) (map[string]int64, error)
	IncrementWithExpiry(ctx context.Context, key string, expiry time.Duration) (int64, error)
	IncrementWithGreaterExpiry(ctx context.Context, key string, value int64, expiry time.Duration) (result int64, err error)

	SetExpiryAll(ctx context.Context, expiry time.Duration, keys ...string) error

	Delete(ctx context.Context, keys ...string) error
	DeleteAll(ctx context.Context) (err error)
	DeleteFromStringSet(ctx context.Context, key string, values ...string) error
	DeleteFromStringList(ctx context.Context, key string, values ...string) error

	Subscribe(ctx context.Context, subscriber Subscriber, options SubscriberOptions, patterns ...string) SubscriberCloser

	Ping(ctx context.Context) error

	GetMultipleHashStringMap(ctx context.Context, keys ...string) (map[string]map[string]string, error)
	SetHashStringMap(ctx context.Context, key string, keyValuePair []string) error
	SetHashStringMapWithGreaterExpiry(ctx context.Context, key string, keyValuePair []string, expiry time.Duration) error
	SetHashStringMapWithGreaterExpiryAndDelete(ctx context.Context, key, deleteKey string, keyValuePair, deleteFields []string, expiry time.Duration) error
	DeleteFromHashStringMap(ctx context.Context, key string, fields ...string) error
	DeleteFromMultipleHashStringMap(ctx context.Context, keyValueHashMap map[string][]string) error

	Close() error
}

// ErrProviderNotSupported is the error used when the provider is not supported
var ErrProviderNotSupported = errors.New("provider not supported")

// ErrNotFound is the error when the data is not found in the cache
var ErrNotFound = errors.New("no data found for the key")

// ErrPastExpiry is the error when the data has expired
var ErrPastExpiry = errors.New("data is past expiry")

// ErrIncorrectValueKind is the error when the key is having the wrong kind of value
var ErrIncorrectValueKind = errors.New("key is holding the wrong kind of value")

// NoExpirationDuration is the value to be used for no expiration
const NoExpirationDuration time.Duration = -1

func getIntConfigOption(options map[string]interface{}, key string) (int, error) {
	var val interface{}
	var ok bool
	var i int
	if options == nil {
		return i, errors.New("no options provided")
	}
	if val, ok = options[key]; ok {
		if i, ok = val.(int); !ok {
			return i, fmt.Errorf("invalid %s, must be an int", key)
		}
	} else {
		return i, fmt.Errorf("missing %s", key)
	}
	return i, nil
}

func getStringConfigOption(options map[string]interface{}, key string) (string, error) {
	var val interface{}
	var ok bool
	var s string
	if options == nil {
		return s, errors.New("no options provided")
	}
	if val, ok = options[key]; ok {
		if s, ok = val.(string); !ok {
			return s, fmt.Errorf("invalid %s, must be a string", key)
		}
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
	return s, nil
}

func getStringSliceConfigOption(options map[string]interface{}, key string) ([]string, error) {
	var val interface{}
	var ok bool
	var s []string
	if options == nil {
		return s, errors.New("no options provided")
	}
	if val, ok = options[key]; ok {
		if s, ok = val.([]string); !ok {
			return s, fmt.Errorf("invalid %s, must be an int", key)
		}
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
	return s, nil
}

func getDurationConfigOption(options map[string]interface{}, key string) (time.Duration, error) {
	var val interface{}
	var ok bool
	var d time.Duration
	if options == nil {
		return d, errors.New("no options provided")
	}
	if val, ok = options[key]; ok {
		if d, ok = val.(time.Duration); !ok {
			return d, fmt.Errorf("invalid %s, must be a duration", key)
		}
	} else {
		return d, fmt.Errorf("missing %s", key)
	}
	return d, nil
}

func (o *SubscriberOptions) init() {
	if o.MaxConcurrency == 0 {
		o.MaxConcurrency = 1
	}
	if o.ProcessTimeout == 0 {
		o.ProcessTimeout = 5 * time.Second
	}
}

// New is used to initialise and get the instance of a config client
func New(ctx context.Context, options Options) (Client, error) {
	if options.Provider == InMemory {
		return newInMemoryCacheClient(ctx, options), nil
	}
	if options.Provider == Redis {
		return newRedisCacheClient(ctx, options)
	}
	if options.Provider == RedisCluster {
		return newRedisClusterCacheClient(ctx, options)
	}
	if options.Provider == LRUInMemory {
		return newLRUInMemoryCacheClient(ctx, options), nil
	}
	return nil, ErrProviderNotSupported
}
