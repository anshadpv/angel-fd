package actuator

// Common Constants
const (
	applicationKey      = "app"
	buildStampKey       = "buildStamp"
	envKey              = "env"
	gitKey              = "git"
	gitCommitAuthorKey  = "commitAuthor"
	gitCommitIDKey      = "commitId"
	gitCommitTimeKey    = "commitTime"
	gitPrimaryBranchKey = "branch"
	gitURLKey           = "url"
	goRoutinesKey       = "goroutine"
	hostNameKey         = "hostName"
	nameKey             = "name"
	usernameKey         = "username"
	versionKey          = "version"
	slash               = "/"
)

// Endpoints
const (
	infoEndpoint       = "/info"
	metricsEndpoint    = "/metrics"
	pingEndpoint       = "/ping"
	threadDumpEndpoint = "/threadDump"
	healthEndpoint     = "/health"
)

// Response constants
const (
	contentTypeHeader          = "Content-Type"
	applicationJSONContentType = "application/json"
	textStringContentType      = "text/string"
)

// Error messages
const (
	methodNotAllowedError = "requested method is not allowed on the called endpoint"
	notFoundError         = "not found"
	profileNotFoundError  = "profile not found"
)
