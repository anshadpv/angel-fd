package v1

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/external"
	"github.com/angel-one/fd-core/factory"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
	upswing external.UpSwing
}

func DefaultTokenController() TokenController {
	return TokenController{upswing: factory.GetUpSwingExternalService()}
}

// GetToken godoc
// @Summary      Get token from vendor
// @Description  Get token from vendor for given clientCode
// @version 1.0
// @Tags         Provider
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Param X-Source header string false "source"
// @Success      200  {object}  model.APIResponse{data=model.PCIRegistrationResponse}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/token/{provider} [GET]
func (c *TokenController) GetToken(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	provider := gctx.Param(constants.Provider)
	log.Info(ctx).Msgf("ClientCode: %s; Provider: %s", clientCode, provider)

	if !slices.Contains(constants.KnownProviders, provider) {
		msg := fmt.Sprintf("Provider %s not supported", provider)
		errors.Throw(gctx, goerr.New(nil, http.StatusForbidden, msg))
		return
	}

	if provider == constants.UpSwingProvider {
		response, err := c.upswing.GetNetWorthData(ctx, clientCode)
		if err != nil {
			errors.Throw(gctx, err)
			return
		}
		log.Trace(ctx).Msgf("Response: %+v", response)
		gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
		return
	}

	errors.Throw(gctx, goerr.New(nil, http.StatusBadRequest, "Unexpected scenario occured"))
}
