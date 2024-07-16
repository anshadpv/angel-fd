package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitRefresherRoute(vGroups ...*gin.RouterGroup) {
	initRefresherV1Group(vGroups[0])
}

func initRefresherV1Group(v1Group *gin.RouterGroup) {
	refreshController := v1.DefaultRefresherController()

	refresh := v1Group.Group(constants.Pending)
	{
		refresh.GET(constants.PathSplitter+constants.PathParam+constants.Refresher, refreshController.GetPendingJourneyRefresher)
	}
}
