package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitUpswingRoute(vGroups ...*gin.RouterGroup) {
	initUpswingV1Group(vGroups[0])
}

func initUpswingV1Group(v1Group *gin.RouterGroup) {
	upswingInfoController := v1.DefaultUpswingController()
	upswng := v1Group.Group(constants.Upswing)
	{
		upswng.GET("", upswingInfoController.GetUpswing)
	}
}
