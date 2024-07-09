package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitHomepageRoute(vGroups ...*gin.RouterGroup) {
	initHomepageV1Group(vGroups[0])
}

func initHomepageV1Group(v1Group *gin.RouterGroup) {
	homeController := v1.DefaultHomepageController()
	plans := v1Group.Group(constants.Home)
	{
		plans.GET("", homeController.GetHomepage)
	}
}
