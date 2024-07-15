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

type UpswingController struct {
	upswing external.UpSwing
}

func DefaultUpswingController() UpswingController {
	return UpswingController{upswing: factory.GetUpSwingExternalService()}
}

func (c *UpswingController) GetUpswing(gctx *gin.Context) {
	var response model.CombinedResponse
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	provider := constants.UpSwingProvider //hardcode for now until we find another vendor
	log.Info(ctx).Msgf("ClientCode: %s; Provider: %s", clientCode, provider)

	if !slices.Contains(constants.KnownProviders, provider) {
		msg := fmt.Sprintf("Provider %s not supported", provider)
		errors.Throw(gctx, goerr.New(nil, http.StatusForbidden, msg))
		return
	}

	netWorthResponse, err := c.upswing.GetNetWorthData(ctx, clientCode)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}

	pendingJourneyResponse, err := c.upswing.GetPendingJourneyData(ctx, clientCode)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}

	response = model.CombinedResponse{
		NetWorthData:   *netWorthResponse,
		PendingJourney: *pendingJourneyResponse,
	}

	log.Trace(ctx).Msgf("Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
