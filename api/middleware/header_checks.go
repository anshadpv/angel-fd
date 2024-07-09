package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/angel-one/fd-core/commons/log"
	"github.com/gin-gonic/gin"
)

// configurable options for header middleware
type HeadersMiddlewareOptions struct {
	WhitelistedHosts string
	ExcludedURI      []string
	originsMap       map[string]bool
}

func (o *HeadersMiddlewareOptions) init() {
	o.WhitelistedHosts = strings.TrimSpace(o.WhitelistedHosts)
	if o.WhitelistedHosts == "" {
		o.WhitelistedHosts = AllHeaderValue
	} else {
		o.originsMap = make(map[string]bool)
		origins := strings.Split(o.WhitelistedHosts, ",")
		for _, v := range origins {
			o.originsMap[v] = true
		}
	}
}

// explict checks on headers that ensures , no tampering is attempted
func HeaderChecks(o HeadersMiddlewareOptions) gin.HandlerFunc {
	o.init()
	return func(ctx *gin.Context) {
		flag := false
		requestURI := ctx.Request.RequestURI
		requestHost := ctx.Request.Host
		url := ctx.Request.RemoteAddr

		// ignore if the request is for excluded URI
		for _, v := range o.ExcludedURI {
			if strings.EqualFold(v, requestURI) {
				flag = true
				break
			}
		}

		// if not on excluded list, then check valid host
		if !flag {
			for key := range o.originsMap {
				match, _ := regexp.MatchString(key, requestHost)
				if match {
					flag = true
					break
				}
			}
		}
		if !flag && len(requestHost) > 0 {
			log.Info(ctx).Msgf("Ignoring request from  Host: %s; URL : %s; URI: %s", requestHost, url, requestURI)
			ctx.AbortWithStatus(http.StatusPreconditionFailed)
			return
		}
		ctx.Next()
	}
}
