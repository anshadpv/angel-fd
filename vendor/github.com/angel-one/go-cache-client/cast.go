package cache

import (
	"context"
	"encoding/xml"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"reflect"
)

func toInt64(d interface{}) (int64, error) {
	return cast.ToInt64E(d)
}

func toFloat64(d interface{}) (float64, error) {
	return cast.ToFloat64E(d)
}

func toString(d interface{}) (string, error) {
	return cast.ToStringE(d)
}

func toBool(d interface{}) (bool, error) {
	return cast.ToBoolE(d)
}

func toSlice(d interface{}) ([]interface{}, error) {
	var s []interface{}
	switch v := d.(type) {
	case []interface{}:
		return append(s, v...), nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	case string:
		err := jsoniter.Unmarshal([]byte(v), &s)
		return s, err
	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", d, d)
	}
}

func toInt64Slice(d interface{}) ([]int64, error) {
	if d == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", d, d)
	}
	switch v := d.(type) {
	case []int64:
		return v, nil
	case string:
		var i []int64
		err := jsoniter.Unmarshal([]byte(v), &i)
		return i, err
	}
	kind := reflect.TypeOf(d).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(d)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toInt64(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", d, d)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", d, d)
	}
}

func toFloat64Slice(d interface{}) ([]float64, error) {
	if d == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []float64", d, d)
	}
	switch v := d.(type) {
	case []float64:
		return v, nil
	case string:
		var f []float64
		err := jsoniter.Unmarshal([]byte(v), &f)
		return f, err
	}
	kind := reflect.TypeOf(d).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(d)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toFloat64(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []float64", d, d)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", d, d)
	}
}

func toStringSlice(d interface{}) ([]string, error) {
	var a []string
	switch v := d.(type) {
	case []interface{}:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []string:
		return v, nil
	case []int8:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []int:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []int32:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []int64:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []float32:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case []float64:
		for _, u := range v {
			a = append(a, cast.ToString(u))
		}
		return a, nil
	case string:
		err := jsoniter.Unmarshal([]byte(v), &a)
		return a, err
	case []error:
		for _, err := range d.([]error) {
			a = append(a, err.Error())
		}
		return a, nil
	case interface{}:
		str, err := cast.ToStringE(v)
		if err != nil {
			return a, fmt.Errorf("unable to cast %#v of type %T to []string", d, d)
		}
		return []string{str}, nil
	default:
		return a, fmt.Errorf("unable to cast %#v of type %T to []string", d, d)
	}
}

func toBoolSlice(d interface{}) ([]bool, error) {
	if d == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", d, d)
	}
	switch v := d.(type) {
	case []bool:
		return v, nil
	case string:
		var b []bool
		err := jsoniter.Unmarshal([]byte(v), &b)
		return b, err
	}
	kind := reflect.TypeOf(d).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(d)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := toBool(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", d, d)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []bool", d, d)
	}
}

func toMap(d interface{}) (map[string]interface{}, error) {
	return cast.ToStringMapE(d)
}

func toInt64Map(d interface{}) (map[string]int64, error) {
	return cast.ToStringMapInt64E(d)
}

func toFloat64Map(d interface{}) (map[string]float64, error) {
	var m = map[string]float64{}
	if d == nil {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]float64", d, d)
	}
	switch v := d.(type) {
	case map[interface{}]interface{}:
		for k, val := range v {
			m[cast.ToString(k)] = cast.ToFloat64(val)
		}
		return m, nil
	case map[string]interface{}:
		for k, val := range v {
			m[k] = cast.ToFloat64(val)
		}
		return m, nil
	case map[string]float64:
		return v, nil
	case string:
		err := jsoniter.Unmarshal([]byte(v), &m)
		return m, err
	}
	if reflect.TypeOf(d).Kind() != reflect.Map {
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", d, d)
	}
	mVal := reflect.ValueOf(m)
	v := reflect.ValueOf(d)
	for _, keyVal := range v.MapKeys() {
		val, err := toFloat64(v.MapIndex(keyVal).Interface())
		if err != nil {
			return m, fmt.Errorf("unable to cast %#v of type %T to map[string]int64", d, d)
		}
		mVal.SetMapIndex(keyVal, reflect.ValueOf(val))
	}
	return m, nil
}

func toStringMap(d interface{}) (map[string]string, error) {
	return cast.ToStringMapStringE(d)
}

func toBoolMap(d interface{}) (map[string]bool, error) {
	return cast.ToStringMapBoolE(d)
}

func toJsonBytes(d interface{}) ([]byte, error) {
	switch s := d.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	default:
		return jsoniter.Marshal(d)
	}
}

func loadAndGetJsonBytes(ctx context.Context, key string, loader Loader) (interface{}, error) {
	d, err := loader(ctx, key)
	if err != nil {
		return nil, err
	}
	return toJsonBytes(d)
}

func castFromJson(b interface{}, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot cast to %#v as it is not a reference", v)
	}
	switch s := b.(type) {
	case string:
		return jsoniter.Unmarshal([]byte(s), v)
	case []byte:
		return jsoniter.Unmarshal(s, v)
	default:
		return fmt.Errorf("unable to cast %#v of type %T to json", b, b)
	}
}

func toXmlBytes(d interface{}) ([]byte, error) {
	switch s := d.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	default:
		return xml.Marshal(d)
	}
}

func loadAndGetXmlBytes(ctx context.Context, key string, loader Loader) (interface{}, error) {
	d, err := loader(ctx, key)
	if err != nil {
		return nil, err
	}
	return toXmlBytes(d)
}

func castFromXml(b interface{}, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot cast to %#v as it is not a reference", v)
	}
	switch s := b.(type) {
	case string:
		return xml.Unmarshal([]byte(s), v)
	case []byte:
		return xml.Unmarshal(s, v)
	default:
		return fmt.Errorf("unable to cast %#v of type %T to xml", b, b)
	}
}

func toYamlBytes(d interface{}) ([]byte, error) {
	switch s := d.(type) {
	case string:
		return []byte(s), nil
	case []byte:
		return s, nil
	default:
		return yaml.Marshal(d)
	}
}

func loadAndGetYamlBytes(ctx context.Context, key string, loader Loader) (interface{}, error) {
	d, err := loader(ctx, key)
	if err != nil {
		return nil, err
	}
	return toYamlBytes(d)
}

func castFromYaml(b interface{}, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot cast to %#v as it is not a reference", v)
	}
	switch s := b.(type) {
	case string:
		return yaml.Unmarshal([]byte(s), v)
	case []byte:
		return yaml.Unmarshal(s, v)
	default:
		return fmt.Errorf("unable to cast %#v of type %T to yaml", b, b)
	}
}

func loadAndGetStringSlice(ctx context.Context, key string, loader Loader) ([]string, error) {
	d, err := loader(ctx, key)
	if err != nil {
		return nil, err
	}
	return toStringSlice(d)
}
