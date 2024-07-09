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

type FAQController struct {
	FAQService service.FAQService
}

func DefaultFAQController() FAQController {
	return FAQController{FAQService: service.DefaultFAQService()}
}

// @Summary      Get FAQs
// @Description  Get all FAQs pertaining to home page or FSI
// @version 1.0
// @Tags         FAQ
// @Produce      json
// @Param Authorization header string true "authorization token"
// @Param X-Request-Id header string true "unique request id"
// @Param tag path string true "home or fsi-name"
// @Success      200  {object}  model.APIResponse{data=model.FAQResponse}
// @Failure	     400  {object}  errors.ErrResponse
// @Failure      500  {object}  errors.ErrResponse
// @Router       /v1/faqs/{tag} [GET]
func (c *FAQController) GetFAQs(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	tag := gctx.Param(constants.Tag)
	log.Debug(ctx).Msgf("ClientCode: %s ", clientCode)

	response, err := c.FAQService.GetFAQDetails(ctx, tag)
	if err != nil {
		errMsg := fmt.Sprintf("unable to get FAQ details due to %v", err)
		errors.Throw(gctx, goerr.New(err, http.StatusInternalServerError, errMsg))
		return
	}
	log.Trace(ctx).Msgf("FAQ Response: %+v", response)
	gctx.JSON(http.StatusOK, model.APIResponse{Data: response})
}
