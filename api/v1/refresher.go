package v1

import (
	"github.com/angel-one/fd-core/business/service"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

type RefresherController struct {
	RefresherService service.RefresherService
}

func DefaultRefresherController() RefresherController {
	return RefresherController{RefresherService: service.DefaultRefresherService()}
}

// Swagger not required as this is internal engg API
func (j *RefresherController) GetPendingJourneyRefresher(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	refresher := gctx.Param(constants.Refresher)
	log.Debug(ctx).Msgf("ClientCode: %s; Refresher: %s ", clientCode, refresher)
	j.RefresherService.InvokePendingJourneyRefresher(ctx, refresher)
}
