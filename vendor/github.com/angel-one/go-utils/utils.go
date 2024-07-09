package utils

import (
	"context"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/angel-one/go-utils/errors"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/angel-one/go-utils/log"
	jsoniter "github.com/json-iterator/go"
)

var numbers = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// GetDataAsBytes is used to get the data as bytes.
func GetDataAsBytes(data io.ReadCloser) ([]byte, error) {
	bytes, err := ioutil.ReadAll(data)
	defer closeData(data)
	return bytes, err
}

// GetDataAsString is used to get the data as string.
func GetDataAsString(data io.ReadCloser) (string, error) {
	bytes, err := GetDataAsBytes(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetJSONData is used to get the JSON data parsed into a struct, make sure you pass the struct by reference.
func GetJSONData(data io.ReadCloser, val interface{}) error {
	bytes, err := GetDataAsBytes(data)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(bytes, val)
}

// MarshalJSON - Marshal returns the JSON encoding of data
func MarshalJSON(data interface{}) ([]byte, error) {
	return jsoniter.Marshal(data)
}

// UnmarshalJSON -- returns the parsed JSON-encoded data in `val`
func UnmarshalJSON(data []byte, val interface{}) error {
	return jsoniter.Unmarshal(data, val)
}

// GetXMLData is used to get the XML data parsed into a struct.
// Make sure you pass the struct by reference.
func GetXMLData(data io.ReadCloser, val interface{}) error {
	bytes, err := GetDataAsBytes(data)
	if err != nil {
		return err
	}
	return xml.Unmarshal(bytes, val)
}

// MarshalXML - returns the XML encoding of v.
func MarshalXML(data interface{}) ([]byte, error) {
	return xml.Marshal(data)
}

// UnmarshalXML - returns the parsed XML-encoded data in `val`
func UnmarshalXML(data []byte, val interface{}) error {
	return xml.Unmarshal(data, val)
}

// GetYMLData is used to get the yml data parsed into a struct.
// Make sure you pass the struct by reference.
func GetYMLData(data io.ReadCloser, val interface{}) error {
	bytes, err := GetDataAsBytes(data)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, val)
}

// MarshalYML - returns the YML encoding of v.
func MarshalYML(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}

// UnmarshalYML - returns the parsed yml-encoded data in `val`
func UnmarshalYML(data []byte, val interface{}) error {
	return yaml.Unmarshal(data, val)
}

// GenerateUUID is used to generate uuid
func GenerateUUID() string {
	return uuid.New().String()
}

// GetRandomNumber is used to generate a random number of given size
func GetRandomNumber(size int) (string, error) {
	b := make([]byte, size)
	n, err := io.ReadAtLeast(rand.Reader, b, size)
	if err != nil {
		return "", err
	}
	if n != size {
		return "", errors.New("unable to generate a random number")
	}
	for i := 0; i < len(b); i++ {
		b[i] = numbers[int(b[i])%len(numbers)]
	}
	return string(b), nil
}

// GetStringVariableFromContext is to get string variable from the context
func GetStringVariableFromContext(ctx context.Context, key string) (string, error) {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			if s == "" {
				return s, errors.New(fmt.Sprintf("missing variable value for %s", key))
			}
			return s, nil
		}
		return "", errors.New(fmt.Sprintf("invalid variable value for %s, found %v", key, v))
	}
	return "", errors.New(fmt.Sprintf("missing variable %s", key))
}

// GetStringVariableFromContextD is to get string variable from context ignoring the error
func GetStringVariableFromContextD(ctx context.Context, key string) string {
	s, _ := GetStringVariableFromContext(ctx, key)
	return s
}

// GetBoolVariableFromContext is to get bool variable from the context
func GetBoolVariableFromContext(ctx context.Context, key string) (bool, error) {
	if v := ctx.Value(key); v != nil {
		if b, ok := v.(bool); ok {
			return b, nil
		}
		return false, errors.New(fmt.Sprintf("invalid variable value for %s, found %v", key, v))
	}
	return false, errors.New(fmt.Sprintf("missing variable %s", key))
}

// GetBoolVariableFromContextD is to get string variable from context ignoring the error
func GetBoolVariableFromContextD(ctx context.Context, key string) bool {
	b, _ := GetBoolVariableFromContext(ctx, key)
	return b
}

func closeData(data io.ReadCloser) {
	err := data.Close()
	if err != nil {
		log.Error(nil).Err(err).Msg("error closing data")
	}
}

// DoesStringOptionsContain is used to check if the value exists in the list of options
func DoesStringOptionsContain(options []string, value string) bool {
	matched := false
	source := cleanString(value)
	for _, s := range options {
		if cleanString(s) == source {
			matched = true
			break
		}
	}
	return matched
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// general copy function to copy all k-v into destination map; use only for small maps
func Copy(dest, src interface{}) {
	dv, sv := reflect.ValueOf(dest), reflect.ValueOf(src)

	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}
