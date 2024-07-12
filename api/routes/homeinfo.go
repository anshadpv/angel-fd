package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitHomeInfoRoute(vGroups ...*gin.RouterGroup) {
	initHomeInfoV1Group(vGroups[0])
}

func initHomeInfoV1Group(v1Group *gin.RouterGroup) {
	homeInfoController := v1.DefaultHomeInfoController()
	plans := v1Group.Group(constants.HomeInfo)

	{
		plans.GET("", homeInfoController.GetHomeInfopage)
	}
}
