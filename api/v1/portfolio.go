package v1

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/angel-one/fd-core/business/model"
	v1 "github.com/angel-one/fd-core/business/service/v1"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/external"
	"github.com/angel-one/fd-core/factory"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type PortfolioController struct {
	upswing   external.UpSwing
	portfolio v1.PortfolioService
}

func DefaultPortfolioController() PortfolioController {
	return PortfolioController{upswing: factory.GetUpSwingExternalService(), portfolio: factory.GetPortfolioService()}
}

// Swagger not required - this API would be decommissioned soon
func (p *PortfolioController) GetNetworth(gctx *gin.Context) {
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
		response, err := p.upswing.GetNetWorthData(ctx, clientCode)
		if err != nil {
			errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, "unable to get networth data from upswing"))
			return
		}
		if response == nil {
			//todo:
		} else {
			log.Trace(ctx).Msgf("NetWorth Response: %+v", response)
			gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
		}
		return
	}

	errors.Throw(gctx, goerr.New(nil, http.StatusBadRequest, "Unexpected scenario occured"))
}

// @Summary      Get Portfolio Summary
// @Description  Get Portfolio summary
// @version 1.0
// @Tags         Portfolio
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Success      200  {object}  model.APIResponse{data=model.NetWorthResponse}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/portfolio [GET]
func (p *PortfolioController) GetPortfolio(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	provider := constants.UpSwingProvider //hardcode for now until we find another vendor
	log.Info(ctx).Msgf("ClientCode: %s; Provider: %s", clientCode, provider)

	if !slices.Contains(constants.KnownProviders, provider) {
		msg := fmt.Sprintf("Provider %s not supported", provider)
		errors.Throw(gctx, goerr.New(nil, http.StatusForbidden, msg))
		return
	}

	response, err := p.portfolio.GetPortfolio(ctx, clientCode, provider)
	if err != nil {
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, "unable to get portfolio data"))
		return
	}
	if response == nil {
		//todo:
	} else {
		log.Trace(ctx).Msgf("Portfolio Response: %+v", response)
		gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
	}
	return

}
