package flags

import (
	"fmt"
	"os"
	"strings"

	"github.com/angel-one/fd-core/constants"
	flag "github.com/spf13/pflag"
)

var (
	mode           = flag.String(constants.ModeKey, constants.ModeDefaultValue, constants.ModeUsage)
	port           = flag.Int(constants.PortKey, constants.PortDefaultValue, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.BaseConfigPathDefaultValue,
		constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

// Env is the application.yml runtime environment
func Env() string {
	env := os.Getenv(constants.EnvKey)
	if env == "" {
		return constants.EnvDefaultValue
	}
	return env
}

// Port is the application.yml port number where the process will be started
func Port() int {
	return *port
}

// BaseConfigPath is the path that holds the configuration files
func BaseConfigPath() string {
	return *baseConfigPath
}

// Mode is the run mode, can be test or release
func Mode() string {
	return *mode
}

// AWSRegion is the region where the application is running
func AWSRegion() string {
	region := os.Getenv(constants.AWSRegionKey)
	if region == "" {
		return constants.AWSRegionDefaultValue
	}
	return region
}

func GetAWSSecretName() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", constants.AWSSecretsName, Env()))
}
