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

type HomepageController struct {
	HomepageService service.HomepageService
}

func DefaultHomepageController() HomepageController {
	return HomepageController{HomepageService: service.DefaultHomepageService()}
}

// @Summary      Get Homepage data
// @Description  Get all data pertaining to home page
// @version 1.0
// @Tags         Home
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Success      200  {object}  model.APIResponse{data=model.Homepage}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/home [GET]
func (c *HomepageController) GetHomepage(gctx *gin.Context) {
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

	response, err := c.HomepageService.GetHomePageDetails(ctx, clientCode, provider)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get homepage details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("Homepage Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
