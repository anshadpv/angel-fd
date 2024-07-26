package cache

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type inMemoryCacheClientData struct {
	data       interface{}
	expiration int64
}

type inMemoryCacheClient struct {
	mu      sync.RWMutex
	signal  chan struct{}
	data    map[string]inMemoryCacheClientData
	hData   map[string]map[string]inMemoryCacheClientData
	loader  Loader
	logger  Logger
	options inMemoryCacheClientOptions
}

type inMemoryCacheClientOptions struct {
	defaultExpiration            time.Duration
	defaultEvictionCheckInterval time.Duration
}

// by default, the eviction check will be done every 30 mins
// it can be overridden by a configuration obviously
var defaultEvictionCheckInterval = 30 * time.Minute

func getInMemoryCacheClientOptions(options map[string]interface{}) inMemoryCacheClientOptions {
	var clientOptions inMemoryCacheClientOptions
	var err error
	clientOptions.defaultExpiration, err = getDurationConfigOption(options, "defaultExpiration")
	if err != nil {
		clientOptions.defaultExpiration = NoExpirationDuration
	}
	clientOptions.defaultEvictionCheckInterval, err = getDurationConfigOption(options, "defaultEvictionCheckInterval")
	if err != nil {
		clientOptions.defaultEvictionCheckInterval = defaultEvictionCheckInterval
	}
	return clientOptions
}

func (i *inMemoryCacheClient) log(ctx context.Context, startTS time.Time, err error) {
	if i.logger != nil {
		i.logger(ctx, getRequiredCallerFunctionName(), time.Now().Sub(startTS).Milliseconds(), err)
	}
}

func (i *inMemoryCacheClient) get(ctx context.Context, key string) (result interface{}, err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.RLock()
	defer i.mu.RUnlock()
	if data, ok := i.data[key]; ok {
		if data.expiration > 0 && time.Now().UnixNano() > data.expiration {
			return nil, ErrPastExpiry
		}
		return data.data, nil
	}
	return nil, ErrNotFound
}

func (i *inMemoryCacheClient) watchForEviction() {
	ticker := time.NewTicker(i.options.defaultEvictionCheckInterval)
	for {
		select {
		case <-ticker.C:
			// try to evict
			i.mu.Lock()
			for k, v := range i.data {
				if v.expiration > 0 && time.Now().UnixNano() > v.expiration {
					delete(i.data, k)
				}
			}
			i.mu.Unlock()
		case <-i.signal:
			// close the watch
			return
		}
	}
}

func newInMemoryCacheClient(_ context.Context, options Options) *inMemoryCacheClient {
	clientOptions := getInMemoryCacheClientOptions(options.Params)
	i := &inMemoryCacheClient{
		mu:      sync.RWMutex{},
		signal:  make(chan struct{}, 1),
		data:    make(map[string]inMemoryCacheClientData),
		loader:  options.Loader,
		logger:  options.Logger,
		options: clientOptions,
	}
	go i.watchForEviction()
	return i
}

func (i *inMemoryCacheClient) Get(ctx context.Context, key string) (interface{}, error) {
	d, err := i.get(ctx, key)
	if err != nil {
		// this means either the data is past expiry, or the data does not exist
		if i.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = i.loader(ctx, key)
			if err != nil {
				return nil, err
			}
			// otherwise, set the data in the cache, also use it
			_ = i.SetD(ctx, key, d)
			return d, nil
		}
		// otherwise, use the error above
		return nil, err
	}
	// we have the value already
	return d, nil
}

func (i *inMemoryCacheClient) GetInt(ctx context.Context, key string) (int64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toInt64(d)
}

func (i *inMemoryCacheClient) GetFloat(ctx context.Context, key string) (float64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return toFloat64(d)
}

func (i *inMemoryCacheClient) GetString(ctx context.Context, key string) (string, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return toString(d)
}

func (i *inMemoryCacheClient) GetBool(ctx context.Context, key string) (bool, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return toBool(d)
}

func (i *inMemoryCacheClient) GetSlice(ctx context.Context, key string) ([]interface{}, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toSlice(d)
}

func (i *inMemoryCacheClient) GetIntSlice(ctx context.Context, key string) ([]int64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Slice(d)
}

func (i *inMemoryCacheClient) GetFloatSlice(ctx context.Context, key string) ([]float64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Slice(d)
}

func (i *inMemoryCacheClient) GetStringSlice(ctx context.Context, key string) ([]string, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringSlice(d)
}

func (i *inMemoryCacheClient) GetBoolSlice(ctx context.Context, key string) ([]bool, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolSlice(d)
}

func (i *inMemoryCacheClient) GetMap(ctx context.Context, key string) (map[string]interface{}, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toMap(d)
}

func (i *inMemoryCacheClient) GetIntMap(ctx context.Context, key string) (map[string]int64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toInt64Map(d)
}

func (i *inMemoryCacheClient) GetFloatMap(ctx context.Context, key string) (map[string]float64, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toFloat64Map(d)
}

func (i *inMemoryCacheClient) GetStringMap(ctx context.Context, key string) (map[string]string, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringMap(d)
}

func (i *inMemoryCacheClient) GetBoolMap(ctx context.Context, key string) (map[string]bool, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return toBoolMap(d)
}

func (i *inMemoryCacheClient) GetAll(ctx context.Context, keys ...string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range keys {
		d, err := i.Get(ctx, key)
		if err == nil {
			result[key] = d
		}
	}
	return result
}

func (i *inMemoryCacheClient) GetJson(ctx context.Context, key string, value interface{}) error {
	d, err := i.get(ctx, key)
	if err != nil {
		// this means either the data is past expiry, or the data does not exist
		if i.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetJsonBytes(ctx, key, i.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = i.SetD(ctx, key, d)
			return castFromJson(d, value)
		}
		// otherwise, use the error above
		return err
	}
	// we have the value already
	return castFromJson(d, value)
}

func (i *inMemoryCacheClient) GetXml(ctx context.Context, key string, value interface{}) error {
	d, err := i.get(ctx, key)
	if err != nil {
		// this means either the data is past expiry, or the data does not exist
		if i.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetXmlBytes(ctx, key, i.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = i.SetD(ctx, key, d)
			return castFromXml(d, value)
		}
		// otherwise, use the error above
		return err
	}
	// we have the value already
	return castFromXml(d, value)
}

func (i *inMemoryCacheClient) GetYaml(ctx context.Context, key string, value interface{}) error {
	d, err := i.get(ctx, key)
	if err != nil {
		// this means either the data is past expiry, or the data does not exist
		if i.loader != nil {
			// now that the loader does not exist, so in that case try to load it
			d, err = loadAndGetYamlBytes(ctx, key, i.loader)
			if err != nil {
				return err
			}
			// otherwise, set the data in the cache, also use it
			_ = i.SetD(ctx, key, d)
			return castFromYaml(d, value)
		}
		// otherwise, use the error above
		return err
	}
	// we have the value already
	return castFromYaml(d, value)
}

func (i *inMemoryCacheClient) GetStringSet(ctx context.Context, key string) ([]string, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v, ok := d.(map[string]struct{}); ok {
		values := make([]string, 0, len(v))
		for k := range v {
			values = append(values, k)
		}
		return values, nil
	} else {
		return nil, ErrIncorrectValueKind
	}
}

func (i *inMemoryCacheClient) GetStringList(ctx context.Context, key string, offset, limit int64) ([]string, error) {
	d, err := i.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v, ok := d.([]string); ok {
		return v[offset:(offset + limit)], nil
	} else {
		return nil, ErrIncorrectValueKind
	}
}

func (i *inMemoryCacheClient) GetMultipleStringLists(ctx context.Context, keys ...string) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, key := range keys {
		d, err := i.Get(ctx, key)
		if err == nil {
			if v, ok := d.([]string); ok {
				result[key] = v
			}
		}
	}
	return result, nil
}

func (i *inMemoryCacheClient) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	i.data[key] = inMemoryCacheClientData{
		data:       value,
		expiration: e,
	}
	return nil
}

func (i *inMemoryCacheClient) SetD(ctx context.Context, key string, value interface{}) error {
	return i.Set(ctx, key, value, i.options.defaultExpiration)
}

func (i *inMemoryCacheClient) SetAll(ctx context.Context, expiry time.Duration, values map[string]interface{}) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	for key, value := range values {
		i.data[key] = inMemoryCacheClientData{
			data:       value,
			expiration: e,
		}
	}
	return nil
}

func (i *inMemoryCacheClient) SetAllD(ctx context.Context, values map[string]interface{}) error {
	return i.SetAll(ctx, i.options.defaultExpiration, values)
}

func (i *inMemoryCacheClient) SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toJsonBytes(value)
	if err != nil {
		return err
	}
	return i.Set(ctx, key, b, expiry)
}

func (i *inMemoryCacheClient) SetJsonD(ctx context.Context, key string, value interface{}) error {
	return i.SetJson(ctx, key, value, i.options.defaultExpiration)
}

func (i *inMemoryCacheClient) SetXml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toXmlBytes(value)
	if err != nil {
		return err
	}
	return i.Set(ctx, key, b, expiry)
}

func (i *inMemoryCacheClient) SetXmlD(ctx context.Context, key string, value interface{}) error {
	return i.SetXml(ctx, key, value, i.options.defaultExpiration)
}

func (i *inMemoryCacheClient) SetYaml(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	b, err := toYamlBytes(value)
	if err != nil {
		return err
	}
	return i.Set(ctx, key, b, expiry)
}

func (i *inMemoryCacheClient) SetYamlD(ctx context.Context, key string, value interface{}) error {
	return i.SetYaml(ctx, key, value, i.options.defaultExpiration)
}

func (i *inMemoryCacheClient) AddToStringSet(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	if d, ok := i.data[key]; ok {
		if d.expiration > e {
			d.expiration = e
		}
		var v map[string]struct{}
		if v, ok = d.data.(map[string]struct{}); ok {
		} else {
			return ErrIncorrectValueKind
		}
		for _, value := range values {
			v[value] = struct{}{}
		}
	} else {
		v := make(map[string]struct{})
		for _, value := range values {
			v[value] = struct{}{}
		}
		i.data[key] = inMemoryCacheClientData{
			data:       v,
			expiration: e,
		}
	}
	return nil
}

func (i *inMemoryCacheClient) AddToStringSetD(ctx context.Context, key string, values ...string) error {
	return i.AddToStringSet(ctx, key, i.options.defaultExpiration, values...)
}

func (i *inMemoryCacheClient) AddToStringList(ctx context.Context, key string, expiry time.Duration, values ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	if d, ok := i.data[key]; ok {
		if d.expiration > e {
			d.expiration = e
		}
		var v []string
		if v, ok = d.data.([]string); ok {
		} else {
			return ErrIncorrectValueKind
		}
		v = append(v, values...)
	} else {
		v := make([]string, 0)
		v = append(v, values...)
		i.data[key] = inMemoryCacheClientData{
			data:       v,
			expiration: e,
		}
	}
	return nil
}

func (i *inMemoryCacheClient) AddToStringListD(ctx context.Context, key string, values ...string) error {
	return i.AddToStringList(ctx, key, i.options.defaultExpiration, values...)
}

func (i *inMemoryCacheClient) IncrementAll(ctx context.Context, expiry time.Duration, keys ...string) (result map[string]int64, err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	data := make(map[string]int64)
	for _, key := range keys {
		if d, ok := i.data[key]; ok {
			if d.expiration > time.Now().UnixNano() {
				v, err := cast.ToInt64E(d.data)
				if err != nil {
					return nil, ErrIncorrectValueKind
				}
				i.data[key] = inMemoryCacheClientData{
					data:       v + 1,
					expiration: d.expiration,
				}
				data[key] = v + 1
				continue
			}
		}
		// data past expiry or missing
		i.data[key] = inMemoryCacheClientData{
			data:       1,
			expiration: time.Now().Add(expiry).UnixNano(),
		}
		data[key] = 1
	}
	return data, nil
}

func (i *inMemoryCacheClient) IncrementIntAllD(ctx context.Context, keys ...string) (map[string]int64, error) {
	return i.IncrementAll(ctx, i.options.defaultExpiration, keys...)
}

func (i *inMemoryCacheClient) IncrementWithExpiry(_ context.Context, _ string,
	_ time.Duration) (int64, error) {
	return 0, nil
}

func (i *inMemoryCacheClient) IncrementWithGreaterExpiry(ctx context.Context, key string, value int64,
	expiry time.Duration) (int64, error) {
	var err error
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()

	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}

	var dataExpiry int64
	if d, ok := i.data[key]; ok {
		dataExpiry = d.expiration
	}

	// take the greater expiry from data.expiration, expiry
	if dataExpiry > e {
		e = dataExpiry
	}

	// increment value by existing data
	if d, ok := i.data[key]; ok {
		val, ok := d.data.(int64)
		if !ok {
			return 0, ErrNotFound
		}
		value += val
	}

	i.data[key] = inMemoryCacheClientData{
		data:       value,
		expiration: e,
	}

	return 0, nil
}

func (i *inMemoryCacheClient) SetExpiryAll(ctx context.Context, expiry time.Duration, keys ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	for _, key := range keys {
		if d, ok := i.data[key]; ok {
			i.data[key] = inMemoryCacheClientData{
				data:       d.data,
				expiration: e,
			}
		} else {
			return ErrNotFound
		}
	}
	return nil
}

func (i *inMemoryCacheClient) Delete(ctx context.Context, keys ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, k := range keys {
		delete(i.data, k)
	}
	return nil
}
func (i *inMemoryCacheClient) DeleteAll(ctx context.Context) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	for key := range i.data {
		delete(i.data, key)
	}
	return nil
}

func (i *inMemoryCacheClient) DeleteFromStringSet(ctx context.Context, key string, values ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	if d, ok := i.data[key]; ok {
		var v map[string]struct{}
		if v, ok = d.data.(map[string]struct{}); ok {
			for _, k := range values {
				delete(v, k)
			}
			return nil
		} else {
			return ErrIncorrectValueKind
		}
	} else {
		return ErrNotFound
	}
}

func (i *inMemoryCacheClient) DeleteFromStringList(ctx context.Context, key string, values ...string) (err error) {
	ts := time.Now()
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	if d, ok := i.data[key]; ok {
		var v []string
		if v, ok = d.data.([]string); ok {
			newV := make([]string, 0, len(v))
			for _, vd := range v {
				found := false
				for _, value := range values {
					if value == vd {
						found = true
						break
					}
				}
				if !found {
					newV = append(newV, vd)
				}
			}
			i.data[key] = inMemoryCacheClientData{
				data:       newV,
				expiration: i.data[key].expiration,
			}
			return nil
		} else {
			return ErrIncorrectValueKind
		}
	} else {
		return ErrNotFound
	}
}

func (i *inMemoryCacheClient) Subscribe(_ context.Context, _ Subscriber, _ SubscriberOptions, _ ...string) SubscriberCloser {
	return nil
}

func (i *inMemoryCacheClient) Ping(ctx context.Context) error {
	return nil
}

func (i *inMemoryCacheClient) Close() error {
	i.signal <- struct{}{}
	i.mu.Lock()
	defer i.mu.Unlock()
	for k := range i.data {
		delete(i.data, k)
	}
	return nil
}

func (i *inMemoryCacheClient) GetMultipleHashStringMap(_ context.Context, keys ...string) (map[string]map[string]string, error) {
	var result map[string]map[string]string
	for _, k := range keys {
		var mp map[string]string
		for d, v := range i.hData[k] {
			if v.data != nil {
				mp[d], _ = v.data.(string)
			}
		}
		result[k] = mp
	}
	return result, nil
}

func (i *inMemoryCacheClient) SetHashStringMap(ctx context.Context, key string, values []string) error {
	ts := time.Now()
	var err error
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var input map[string]inMemoryCacheClientData
	for idx := 0; idx < len(values)-1; idx += 2 {
		input[values[idx]] = inMemoryCacheClientData{
			data: values[idx+1],
		}
	}
	i.hData[key] = input
	return nil
}

func (i *inMemoryCacheClient) SetHashStringMapWithGreaterExpiry(ctx context.Context, key string, values []string, expiry time.Duration) error {
	ts := time.Now()
	var err error
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	var input map[string]inMemoryCacheClientData
	for idx := 0; idx < len(values)-1; idx += 2 {
		input[values[idx]] = inMemoryCacheClientData{
			data:       values[idx+1],
			expiration: e,
		}
	}
	i.hData[key] = input
	return nil
}

func (i *inMemoryCacheClient) DeleteFromHashStringMap(ctx context.Context, key string, fields ...string) error {
	ts := time.Now()
	var err error
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	if d, ok := i.hData[key]; ok {
		if d == nil {
			return ErrIncorrectValueKind
		}
		for k, _ := range d {
			for _, f := range fields {
				if k == f {
					delete(d, f)
				}
			}
		}
		i.hData[key] = d
	} else {
		return ErrNotFound
	}
	return nil
}

func (i *inMemoryCacheClient) SetHashStringMapWithGreaterExpiryAndDelete(ctx context.Context, key, deleteKey string, values, deleteFields []string, expiry time.Duration) error {
	ts := time.Now()
	var err error
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	var e int64
	if expiry > 0 {
		e = time.Now().Add(expiry).UnixNano()
	}
	var input map[string]inMemoryCacheClientData
	for idx := 0; idx < len(values)-1; idx += 2 {
		input[values[idx]] = inMemoryCacheClientData{
			data:       values[idx+1],
			expiration: e,
		}
	}
	i.hData[key] = input
	if d, ok := i.hData[deleteKey]; ok {
		if d == nil {
			return ErrIncorrectValueKind
		}
		for k, _ := range d {
			for _, f := range deleteFields {
				if k == f {
					delete(d, f)
				}
			}
		}
		i.hData[deleteKey] = d
	} else {
		return ErrNotFound
	}
	return nil
}

func (i *inMemoryCacheClient) DeleteFromMultipleHashStringMap(ctx context.Context, keyValueHashMap map[string][]string) error {
	ts := time.Now()
	var err error
	defer i.log(ctx, ts, err)
	i.mu.Lock()
	defer i.mu.Unlock()
	for key, val := range keyValueHashMap {
		if d, ok := i.hData[key]; ok {
			if d == nil {
				return ErrIncorrectValueKind
			}
			for k, _ := range d {
				for _, f := range val {
					if k == f {
						delete(d, f)
					}
				}
			}
			i.hData[key] = d
		} else {
			return ErrNotFound
		}
	}

	return nil
}
