package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitPlansRoute(vGroups ...*gin.RouterGroup) {
	initPlansV1Group(vGroups[0])
}

func initPlansV1Group(v1Group *gin.RouterGroup) {
	plansController := v1.DefaultPlansController()

	plans := v1Group.Group(constants.Plans)
	{
		plans.GET("", plansController.GetPlans)
		plans.GET(constants.PathParam+constants.FSI, plansController.GetFSIPlans)
	}
}
