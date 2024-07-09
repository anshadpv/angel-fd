package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	XFrameOptionsHeader           = "X-Frame-Options"
	XFrameOptionsHeaderValue      = "SAMEORIGIN"
	ReferrerPolicyHeader          = "Referrer-Policy"
	ReferrerPolicyHeaderValue     = "strict-origin"
	ContentSecurityHeader         = "Content-Security-Policy"
	ContentSecurityHeaderValue    = "frame-ancestors 'self' https://*.angelone.in;"
	StrictTransportSecurityHeader = "Strict-Transport-Security"
	XContentTypeOptionsHeader     = "X-Content-Type-Options"
	StrictTransportSecurityValue  = "max-age=31536000"
	XContentTypeOptionsValue      = "nosniff"
)

func SecurityHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header(XFrameOptionsHeader, XFrameOptionsHeaderValue)
		ctx.Header(ReferrerPolicyHeader, ReferrerPolicyHeaderValue)
		ctx.Header(ContentSecurityHeader, ContentSecurityHeaderValue)
		ctx.Header(StrictTransportSecurityHeader, StrictTransportSecurityValue)
		ctx.Header(XContentTypeOptionsHeader, XContentTypeOptionsValue)
		ctx.Next()
	}
}
