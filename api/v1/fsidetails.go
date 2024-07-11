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

type FsiDetailsController struct {
	FsiDetailService service.FsiDetailService
}

func DefaultFsiDetailsController() FsiDetailsController {
	return FsiDetailsController{FsiDetailService: service.DefaultFsiDetailService()}
}

func (c *FsiDetailsController) GetFsiDetails(gctx *gin.Context) {
	fmt.Println("REACHED HERE !!!!")
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	queryParams := gctx.Request.URL.Query()
	if len(queryParams) == 0 {
		errors.Throw(gctx, goerr.New(nil, http.StatusBadRequest, "No query parameters provided"))
		return
	}
	var fsiDetailsKeys []string
	var fsiDetailsValues []string

	for key, val := range queryParams {
		fsiDetailsKeys = append(fsiDetailsKeys, key)
		for _, v := range val {
			fsiDetailsValues = append(fsiDetailsValues, v)
		}
	}

	response, err := c.FsiDetailService.GetFsiDetails(ctx, fsiDetailsKeys, fsiDetailsValues)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get compare fsi details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("Compare FSI Details Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
