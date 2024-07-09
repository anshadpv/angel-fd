package middleware

import (
	c "context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angel-one/fd-core/commons/context"
	fderr "github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/errors"
	logConstants "github.com/angel-one/go-utils/constants"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(signingKey []byte) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		if isExcludedPath(gctx) {
			gctx.Next()
			return
		}
		fdctx := context.Build(gctx)

		token := getTokenFromHeader(gctx.Request.Header.Get(constants.HeaderAuthorization))
		if token == constants.Empty {
			fderr.Throw(gctx, errors.HeaderAuthMissingInvalid)
			gctx.Abort()
			return
		}

		if token == "" {
			log.Error(fdctx).Msg("token is empty")
			logHeaders(fdctx, gctx)
			fderr.Throw(gctx, errors.HeaderAuthMissingInvalid)
			return
		}

		claims, err := validateToken(token, signingKey)
		if err != nil {
			log.Error(fdctx).Err(err).Stack().Msg("token validation failed")
			logHeaders(fdctx, gctx)
			fderr.Throw(gctx, errors.NotAuthorized)
			return
		}

		userData, ok := claims[constants.AuthJWTClaimsUserData].(map[string]interface{})
		if !ok {
			log.Error(fdctx).Msg("token claim is not as expected")
			logHeaders(fdctx, gctx)
			fderr.Throw(gctx, errors.NotAuthorized)
			return
		}

		userID, ok := userData[constants.AuthJWTClaimsUserDataUserID].(string)
		if !ok {
			log.Error(fdctx).Msg("token claim-userdata is not as expected")
			logHeaders(fdctx, gctx)
			fderr.Throw(gctx, errors.NotAuthorized)
			return
		}

		if userID == "" {
			userID = constants.AuthGuestUserID
		}

		if userID == constants.AuthGuestUserID {
			log.Error(fdctx).Msg("guest user is not allowed")
			logHeaders(fdctx, gctx)
			fderr.Throw(gctx, errors.NotAuthorized)
			return
		}

		gctx.Request.Header.Add(constants.HeaderUserID, userID)
		gctx.Next()
	}
}

func getTokenFromHeader(bearerToken string) string {
	var accessToken string

	tokenString := strings.Split(bearerToken, " ")
	if len(tokenString) == 2 && tokenString[0] == constants.HeaderAuthorizationBearer {
		accessToken = tokenString[1]
	}

	return accessToken
}

func validateToken(authToken string, signingKey []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %v",
				errors.AuthInvalidSigningMethod,
				token.Header[constants.HeaderAuthorizationAlgorithm])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.NotAuthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.NotAuthorized
}

func logHeaders(sgbctx c.Context, gctx *gin.Context) {
	if reqHeadersBytes, err := json.Marshal(gctx.Request.Header); err != nil {
		log.Info(sgbctx).Msg("Could not Marshal Req Headers")
	} else {
		log.Info(sgbctx).Msgf("Headers: %s", string(reqHeadersBytes))
	}
}

func logAuthErrorDetailed(ctx *gin.Context, err error) {
	path := ctx.Request.URL.Path
	fullPath := ctx.FullPath()
	fullPath = strings.ReplaceAll(fullPath, "/", "")
	fullPath = strings.ReplaceAll(fullPath, ":", "_")
	// log the initial request
	event := log.Info(ctx).
		Str(logConstants.MethodLogParam, ctx.Request.Method).
		Str(logConstants.PathLogParam, path).
		Str(constants.NORMALISED_PATH, fullPath).
		Err(err)
	event.Str(logConstants.QueryLogParam, ctx.Request.URL.RawQuery)
	event.Interface(logConstants.HeaderLogParams, ctx.Request.Header)
	event.Send()
}
