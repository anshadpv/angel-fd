package events

import (
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

type WebhooksController struct {
	webhookService service.WebhookService
}

func DefaultWebhooksController() WebhooksController {
	return WebhooksController{webhookService: service.DefaultWebhookService()}
}

func (c *WebhooksController) ReadUpSwingMessage(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	if clientCode != constants.UpSwingProvider {
		errors.Throw(gctx, goerr.New(nil, http.StatusForbidden, "invalid token"))
		return
	}
	log.Debug(ctx).Msg("New upswing webhook event received")

	var request model.UpSwingWebhookEvent
	if err := gctx.ShouldBind(&request); err != nil {
		errors.Throw(gctx, err)
		return
	}

	log.Debug(ctx).Msgf("Webhook request payload: %+v", request)

	err := c.webhookService.RegisterNewEvent(ctx, constants.UpSwingProvider, request)
	if err != nil {
		errors.Throw(gctx, err)
		return
	}

	gctx.JSON(http.StatusOK, model.EmptyJSON{})
}
