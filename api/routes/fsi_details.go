package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitFsiDetails(vGroups ...*gin.RouterGroup) {
	initFsiDetailsV1Group(vGroups[0])
}

func initFsiDetailsV1Group(v1Group *gin.RouterGroup) {
	fsiDetailsController := v1.DefaultFsiDetailsController()
	fsiDetails := v1Group.Group(constants.Fsi)
	{
		fsiDetails.GET("/details", fsiDetailsController.GetFsiDetails)
	}
}
