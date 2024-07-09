package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitJobsRoute(vGroups ...*gin.RouterGroup) {
	initJobsV1Group(vGroups[0])
}

func initJobsV1Group(v1Group *gin.RouterGroup) {
	jobsController := v1.DefaultJobsController()

	jobs := v1Group.Group(constants.Jobs)
	{
		jobs.GET(constants.PathSplitter+constants.Update+constants.Portfolio+constants.PathParam+constants.Refresher, jobsController.GetPortfolioJob)
		jobs.GET(constants.PathSplitter+constants.Update+constants.PendingJourney+constants.PathParam+constants.Refresher, jobsController.GetPendingJourneyJob)
	}
}
