package constants

// Header Params
const (
	RequestIDHeader                     = "X-requestId"
	AccessControlAllowOriginHeader      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentialsHeader = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersHeader     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsHeader     = "Access-Control-Allow-Methods"
	AllHeaderValue                      = "*"
	TrueHeaderValue                     = "true"
	FalseHeaderValue                    = "false"
	AllMethodsHeaderValue               = "POST,HEAD,PATCH,OPTIONS,GET,PUT"

	// Security
	XFrameOptionsHeader                                                = "X-Frame-Options"
	SameOriginXFrameOptionsHeaderValue                                 = "SAMEORIGIN"
	ReferrerPolicyHeader                                               = "Referrer-Policy"
	StrictOriginReferrerPolicyHeaderValue                              = "strict-origin"
	ContentSecurityPolicyHeader                                        = "Content-Security-Policy"
	AngelOneAllSubDomainFrameAncestorsContentSecurityPolicyHeaderValue = "frame-ancestors 'self' https://*.angelone.in;"
	StrictTransportSecurityHeader                                      = "Strict-Transport-Security"
	StrictTransportSecurityHeaderValue                                 = "max-age=31536000"
	ContentTypeOptionsHeader                                           = "X-Content-Type-Options"
	ContentTypeOptionsHeaderValue                                      = "nosniff"
)
