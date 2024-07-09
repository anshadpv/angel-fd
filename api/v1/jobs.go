package v1

import (
	"github.com/angel-one/fd-core/business/service"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

type JobsController struct {
	JobsService service.JobsService
}

func DefaultJobsController() JobsController {
	return JobsController{JobsService: service.DefaultJobsService()}
}

// Swagger not required as this is internal engg API
func (j *JobsController) GetPortfolioJob(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	refresher := gctx.Param(constants.Refresher)
	log.Debug(ctx).Msgf("ClientCode: %s; Refresher: %s ", clientCode, refresher)

	j.JobsService.InvokePortfolioJob(ctx, refresher)
}

// Swagger not required as this is internal engg API
func (j *JobsController) GetPendingJourneyJob(gctx *gin.Context) {
	ctx := context.Build(gctx)
	clientCode := context.Get(ctx).UserID
	refresher := gctx.Param(constants.Refresher)
	log.Debug(ctx).Msgf("ClientCode: %s; Refresher: %s ", clientCode, refresher)

	j.JobsService.InvokePendingJourneyJob(ctx, refresher)
}
