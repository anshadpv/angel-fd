package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

const (
	AccessControlAllowOriginHeader      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentialsHeader = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersHeader     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsHeader     = "Access-Control-Allow-Methods"
	AllHeaderValue                      = "*"
	TrueHeaderValue                     = "true"
	FalseHeaderValue                    = "false"
	AllMethodsHeaderValue               = "POST,HEAD,PATCH,OPTIONS,GET,PUT"
	OriginRequestHeader                 = "Origin"
)

// CORSMiddlewareOptions is the set of configurable for cors middleware
// this is fork of the common library code
type CORSMiddlewareOptions struct {
	AllowedOrigins   string
	allowCredentials string
	BlockCredentials bool
	AllowedHeaders   string
	AllowedMethods   string
	originsMap       map[string]bool
}

// init is used to set defaults to the options
func (o *CORSMiddlewareOptions) init() {
	o.AllowedOrigins = strings.TrimSpace(o.AllowedOrigins)
	if o.AllowedOrigins == "" {
		o.AllowedOrigins = AllHeaderValue
	} else {
		o.originsMap = make(map[string]bool)
		origins := strings.Split(o.AllowedOrigins, ",")
		for _, v := range origins {
			o.originsMap[v] = true
		}
	}
	if o.BlockCredentials {
		o.allowCredentials = FalseHeaderValue
	} else {
		o.allowCredentials = TrueHeaderValue
	}
	o.AllowedHeaders = strings.TrimSpace(strings.ToUpper(o.AllowedHeaders))
	if o.AllowedHeaders == "" {
		o.AllowedHeaders = AllHeaderValue
	}
	o.AllowedMethods = strings.TrimSpace(strings.ToUpper(o.AllowedMethods))
	if o.AllowedMethods == "" {
		o.AllowedMethods = AllMethodsHeaderValue
	}
}

// CORS is used to allow CORS for the requests this is added to
func CORS(options CORSMiddlewareOptions) gin.HandlerFunc {
	options.init()
	return func(ctx *gin.Context) {
		if strings.EqualFold(options.AllowedOrigins, constants.AllowedOriginsForCorsDefault) {
			ctx.Header(AccessControlAllowOriginHeader, options.AllowedOrigins)
		} else {
			requestOrigin := ctx.Request.Header.Get(OriginRequestHeader)
			if options.originsMap[requestOrigin] {
				ctx.Header(AccessControlAllowOriginHeader, requestOrigin)
			} else { //go by regex matches this is costly operation and should be avoided
				for key := range options.originsMap {
					match, _ := regexp.MatchString(key, requestOrigin)
					if match {
						ctx.Header(AccessControlAllowOriginHeader, requestOrigin)
						break
					}
				}
			}

		}
		ctx.Header(AccessControlAllowCredentialsHeader, options.allowCredentials)
		ctx.Header(AccessControlAllowHeadersHeader, options.AllowedHeaders)
		ctx.Header(AccessControlAllowMethodsHeader, options.AllowedMethods)
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
