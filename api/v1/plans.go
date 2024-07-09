package v1

import (
	"fmt"
	"net/http"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/service"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type PlansController struct {
	PlansService service.PlansService
}

func DefaultPlansController() PlansController {
	return PlansController{PlansService: service.DefaultPlansService()}
}

// @Summary      Get all plans & details
// @Description  Fetch all plan details across all banks/FSIs
// @version 1.0
// @Tags         Plans
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Success      200  {object}  model.APIResponse{data=model.Plans}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/plans [GET]
func (c *PlansController) GetPlans(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	response, err := c.PlansService.GetAllPlans(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get plan details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}

	log.Trace(ctx).Msgf("GetAllPlans Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}

// @Summary      Get FSI plans & details
// @Description  Fetches all plans and details for the specified bank/FSI
// @version 1.0
// @Tags         Plans
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Param fsi path string true "fsi name"
// @Success      200  {object}  model.APIResponse{data=model.FsiPlans}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/plans/{fsi} [GET]
func (c *PlansController) GetFSIPlans(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	fsi := gctx.Param(constants.FSI)
	log.Info(ctx).Msgf("ClientCode: %s; FSI: %s", clientCode, fsi)

	response, err := c.PlansService.GetFSIPlans(ctx, fsi)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get fsi plan details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}

	log.Trace(ctx).Msgf("GetFsiPlans Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
