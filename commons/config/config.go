package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/angel-one/fd-core/commons/flags"
	"github.com/angel-one/fd-core/constants"
	configs "github.com/angel-one/go-config-client"
	"github.com/spf13/cast"
)

type Client struct {
	configs.Client
	secret  string
	env     map[string]string
	Secrets map[string]string
}

var client *Client

func Default() *Client {
	return client
}

func InitTestMode(directory string, configNames ...string) error {
	secretsName := "secrets"
	c, err := configs.New(configs.Options{
		Provider: configs.FileBased,
		Params: map[string]interface{}{
			constants.ConfigDirectoryKey:     directory,
			constants.ConfigNamesKey:         configNames,
			constants.ConfigTypeKey:          constants.ConfigType,
			constants.ConfigSecretsDirectory: directory,
			constants.ConfigSecretNames:      []string{secretsName},
			constants.ConfigSecretType:       constants.ConfigTypeJSON,
		},
	})
	return initClient(c, secretsName, err)
}

func InitReleaseMode(configNames ...string) error {
	secretsName := flags.GetAWSSecretName()
	c, err := configs.New(GetOptions(secretsName, configNames...))
	return initClient(c, secretsName, err)
}

func GetOptions(secretsName string, configNames ...string) configs.Options {
	return configs.Options{
		Provider: configs.AWSAppConfig,
		Params: map[string]interface{}{
			constants.ConfigIDKey:           constants.ApplicationName,
			constants.ConfigAppKey:          constants.ApplicationName,
			constants.ConfigEnvKey:          strings.ToLower(flags.Env()),
			constants.ConfigTypeKey:         constants.ConfigType,
			constants.ConfigNamesKey:        configNames,
			constants.ConfigRegionKey:       flags.AWSRegion(),
			constants.ConfigSecretNames:     []string{secretsName},
			constants.ConfigCredentialsMode: configs.AppConfigSharedCredentialMode,
		},
	}
}

func initClient(c configs.Client, secretName string, err error) error {
	if err != nil {
		return err
	}
	client = getClient(c)
	client.secret = secretName
	client.Secrets = initStore()
	return nil
}

func (c *Client) getStringSecretD(key string, defaultValue string) string {
	val, err := c.GetStringSecret(c.secret, key)
	if err != nil {
		return defaultValue
	}
	return val
}

func (c *Client) GetStringWithEnv(config, key string) (string, error) {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return s, err
	}
	// now time to look for and replace with all the environment variables
	s = c.ReplaceWithEnv(s)
	return s, nil
}

func (c *Client) GetIntFromMap(options map[string]interface{}, key string, defaultv int) (int, error) {
	var val interface{}
	var ok bool
	var s int
	if val, ok = options[key]; ok {
		if s, ok = val.(int); !ok {
			return s, fmt.Errorf("invalid %s, must be a int", key)
		}
	}
	return s, nil
}

func (c *Client) GetStringFromMap(options map[string]interface{}, key string, defaultv string) (string, error) {
	var val interface{}
	var ok bool
	var s string
	if val, ok = options[key]; ok {
		if s, ok = val.(string); !ok {
			return s, fmt.Errorf("invalid %s, must be a string", key)
		}
	} else {
		s = defaultv
	}
	return s, nil
}

func (c *Client) GetStringWithSecretsFromMap(options map[string]interface{}, key string, defaultv string) (string, error) {
	s, err := c.GetStringFromMap(options, key, defaultv)
	return c.ReplaceWithSecret(s), err
}

// Replaces the specified string with the secrets variable values for all name placeholders
func (c *Client) ReplaceWithSecret(value string) string {
	var s string = value
	for k, v := range c.Secrets {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}

func (c *Client) GetStringWithEnvFromMap(options map[string]interface{}, key string, defaultv string) (string, error) {
	s, err := c.GetStringFromMap(options, key, defaultv)
	return c.ReplaceWithEnv(s), err
}

func (i *Client) GetConfigOptionMap(options map[string]interface{}, key string) map[string]string {
	return cast.ToStringMapString(options[key])
}

func (i *Client) GetConfigEnv(options map[string]interface{}, key string) string {
	return i.ReplaceWithEnv(cast.ToString(options[key]))
}

// Replaces the specified string with the environment variable values for all name placeholders
func (c *Client) ReplaceWithEnv(value string) string {
	var s string = value
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}

// Returns present environment name from environment variable
func (c *Client) GetEnv() string {
	return strings.ToLower(flags.Env())
}

// Release true if the current running mode has release argument set
func (c *Client) IsReleaseMode() bool {
	return flags.Mode() == constants.ReleaseMode
}

// Reloads all environment variable once again
func (c *Client) ReloadEnvironment() {
	c.env = getEnvironment()
}

// Return the newly constructed config client reference
func getClient(c configs.Client) *Client {
	return &Client{Client: c, env: getEnvironment()}
}

func getEnvironment() map[string]string {
	env := os.Environ()
	result := make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) == 2 {
			result[s[0]] = s[1]
		}
	}
	return result
}

func GetAllConfigs() map[string]interface{} {
	var dbConfig = make(map[string]interface{})
	dbConfig[constants.DATABASE_NAME] = os.Getenv(constants.DATABASE_NAME)
	dbConfig[constants.DATABASE_USERNAME] = os.Getenv(constants.DATABASE_USERNAME)

	var configMap = make(map[string]interface{})
	configMap["db"] = dbConfig

	return configMap
}

// GetStringWithEnvD is used to get the config with default value by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *Client) GetStringWithEnvD(config, key, defaultValue string) string {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return defaultValue
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}
