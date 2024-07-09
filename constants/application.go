package constants

const ApplicationName = "fd"

// Resource files
const (
	ApplicationConfig = "application"
	LoggerConfig      = "logger"
	DatabaseConfig    = "database"
	HTTPClientConfig  = "http-client"
)

const (
	ConfigsecretNames = "secretNames"
)

// ENV constants
const (
	EnvKey          = "ENV"
	EnvDEV          = "DEV"
	EnvUAT          = "UAT"
	EnvPERF         = "PERF"
	EnvCug          = "CUG"
	EnvProd         = "PROD"
	EnvDefaultValue = "DEV"
)

// Flag constants
const (
	EnvUsage                   = "application.yml runtime environment"
	PortKey                    = "port"
	PortDefaultValue           = 8080
	PortUsage                  = "port"
	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "resources"
	BaseConfigPathUsage        = "path to folder that stores your configurations"
	EnvAuthKey                 = "JWT_SYMMETRIC_KEY"
	ReleaseMode                = "release"
	ModeKey                    = "mode"
	ModeUsage                  = "run mode of the application, can be test or release"
	ModeDefaultValue           = "test"
)

// Error Messages
const (
	ErrAuthKeyEnvNotSet     = "environment key JWT_SYMMETRIC_KEY not set"
	ErrNoUserDetailsFound   = "user not found, please contact service"
	WhitelistedHostIsNotSet = "Whitelisted HOSTS config is missing"
)

const (
	ConfigIDKey            = "id"
	ConfigRegionKey        = "region"
	ConfigAccessKeyID      = "accessKeyId"
	ConfigSecretKey        = "secretKey"
	ConfigCredentialsMode  = "credentialsMode"
	ConfigAppKey           = "app"
	ConfigEnvKey           = "env"
	ConfigTypeKey          = "configType"
	ConfigSecretType       = "secretType"
	ConfigNamesKey         = "configNames"
	ConfigSecretNames      = "secretNames"
	ConfigDirectoryKey     = "configsDirectory"
	ConfigSecretsDirectory = "secretsDirectory"
	ConfigType             = "yaml"
	ConfigTypeJSON         = "json"
	LogLevelKey            = "level"
)

const (
	AWSRegionKey          = "AWS_REGION"
	AWSRegionDefaultValue = "ap-south-1"
	AWSAccessKeyID        = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey    = "AWS_SECRET_ACCESS_KEY"
	AWSBucket             = "AWS_BUCKET"
	AWSSecretsName        = "fd-core"
)

const (
	LogLevelConfigKey = "level"
)

const (
	Empty = ""
)

// http constants
const (
	MethodKey = "method"
	UrlKey    = "url"
)

// http client configs
const (
	UpSwingGenerateToken   = "upswingGenerateToken"
	UpSwingPCIRegistration = "upswingPCIRegistration"
	UpSwingNetWorth        = "upswingNetWorth"
	UpswingDataIngestion   = "upswingDataIngestion"
	UpswingPendingJourney  = "upswingPendingJourney"

	UpPCIField = "{pci}"

	ProfileServerConfig = "profileServiceConfig"
)

const (
	PortfolioUpdateBatchSize = "portfolioUpdateBatchSize"
	PortfolioProvider        = "portfolioProvider"
)

const (
	PendingJourneyUpdateBatchSize = "pendingJourneyUpdateBatchSize"
	PendingJourneyProvider        = "pendingJourneyProvider"
)
