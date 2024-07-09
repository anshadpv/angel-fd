package v1

import (
	"fmt"
	"net/http"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/service"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/errors"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/goerr"
	"github.com/gin-gonic/gin"
)

type CompareController struct {
	CompareService service.CompareService
}

func DefaultCompareController() CompareController {
	return CompareController{CompareService: service.DefaultCompareService()}
}

// @Summary      Get comparable FSIs
// @Description  Fetches all FSIs that are eligible for comparision
// @version 1.0
// @Tags         Compare
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Success      200  {object}  model.APIResponse{data=model.FsiList}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/compare/list [GET]
func (c *CompareController) GetCompareList(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	response, err := c.CompareService.GetCompareList(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get compare list due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("Compare List Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}

// @Summary      Perform comparision
// @Description  Compare FD details for specified FSIs
// @version 1.0
// @Tags         Compare
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Param fsi1 query string true "fsi name"
// @Param fsi2 query string true "fsi name"
// @Success      200  {object}  model.APIResponse{data=map[string]model.CompareFSIDetails}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/compare [GET]
func (c *CompareController) GetCompareDetails(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	queryParams := gctx.Request.URL.Query()
	if len(queryParams) == 0 {
		errors.Throw(gctx, goerr.New(nil, http.StatusBadRequest, "No query parameters provided"))
		return
	}
	var compareFsiKeys []string
	var compareFsiValues []string

	for key, val := range queryParams {
		compareFsiKeys = append(compareFsiKeys, key)
		for _, v := range val {
			compareFsiValues = append(compareFsiValues, v)
		}
	}

	response, err := c.CompareService.GetCompareFsiDetails(ctx, compareFsiKeys, compareFsiValues)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get compare fsi details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("Compare FSI Details Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
