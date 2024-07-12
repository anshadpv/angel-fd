package v1

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/service"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type HomeInfoController struct {
	HomeInfoService service.HomeInfoService
}

func DefaultHomeInfoController() HomeInfoController {
	return HomeInfoController{HomeInfoService: service.DefaultHomeInfoService()}
}

func (c *HomeInfoController) GetHomeInfopage(gctx *gin.Context) {
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

	response, err := c.HomeInfoService.GetHomeInfoDetails(ctx, clientCode, provider)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get homepage details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("Homepage Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
