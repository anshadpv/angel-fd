package constants

import u "github.com/angel-one/go-utils/constants"

// headers
const (
	HeaderAuthorization          = "Authorization"
	HeaderAuthorizationBearer    = "Bearer"
	HeaderToken                  = "token"
	HeaderIPAddress              = "ip-address"
	HeaderAuthtoken              = "Authtoken"
	HeaderAuthorizationAlgorithm = "alg"
	HeaderAppVersion             = "app-version"
	HeaderOsType                 = "os-type"
	HeaderRequestID              = "X-Request-ID"
	HeaderSource                 = "X-Source"
	HeaderClientCode             = "X-Client-Code"
	HeaderPlatform               = "X-Platform"
	HeaderUserID                 = "userID"
	HeaderAppName                = "ApplicationName"
	HeaderDeviceID               = "X-Device-ID"
	HeaderClientIP               = "ClientIP"
	HeaderDeviceType             = "DeviceType"
	HeaderLocation               = "X-Location"
	HeaderMACAddress             = "X-MACAddress"
	HeaderOS                     = "X-OperatingSystem"
	HeaderSourceID               = "X-Source-ID"
	HeaderAcceptLanguage         = "Accept-Language"
)

// URL path constants
const (
	V1            = "/v1"
	V2            = "/v2"
	SwaggerRoute  = "/swagger/*any"
	ActuatorRoute = "/actuator/*any"
	ActuatorInfo  = "/actuator/info"
	Test          = "/test"
	PathSplitter  = "/"
	PathParam     = "/:"
)

const (
	AllowedMethodsForCors         = u.AllMethodsHeaderValue + ",DELETE"
	AllowedOriginsForCorsDefault  = "*"
	AllowedOrginsForCorsConfigKey = "orginsAllowedForCors"
	WhitelistedHostsHeader        = "whitelistedHostHeader"
	WhitelistedHostsHeaderDefault = "angelone.in"
)

// Auth constants
const (
	AuthJWTClaimsUserData       = "userData"
	AuthJWTClaimsUserDataUserID = "user_id"
	AuthGuestUserID             = "guest"
)

const (
	NORMALISED_PATH = "normalisedPath"
)

// URL Paths
const (
	UpSwingWebhookPath = "/external/capture/event"

	Webhook        = "/webhook"
	Token          = "/token"
	Portfolio      = "/portfolio"
	Plans          = "/plans"
	Home           = "/home"
	Networth       = "/networth"
	FAQ            = "/faqs"
	Compare        = "/compare"
	List           = "/list"
	Jobs           = "/jobs"
	Update         = "/update"
	PendingJourney = "/pendingJourney"
	Fsi            = "/fsi"
	Details        = "/details"
	HomeInfo       = "/homeInfo"
	Upswing        = "/upswing"
)

const (
	Provider      = "provider"
	FSI           = "fsi"
	StatusSuccess = "success"
	Tag           = "tag"
	Refresher     = "refresher"
)

var (
	UpSwingProvider = "upswing"
	KnownProviders  = []string{UpSwingProvider}
)

const (
	ErrorCode         = "errorCode"
	ErrClientNotFound = "INTERNAL_CUSTOMER_DETAILS_NOT_FOUND_FOR_PCI"
	ErrPciNotFound    = "[PCI:%s] not found"
)
