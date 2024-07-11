package routes

import (
	"fmt"

	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitFsiDetails(vGroups ...*gin.RouterGroup) {
	fmt.Println("IM HERE")
	initFsiDetailsV1Group(vGroups[0])
}

func initFsiDetailsV1Group(v1Group *gin.RouterGroup) {
	fmt.Println("IM HERE 2")
	fsiDetailsController := v1.DefaultFsiDetailsController()
	fsiDetails := v1Group.Group(constants.Fsi)
	{
		fsiDetails.GET("/details", fsiDetailsController.GetFsiDetails)
	}
}
