package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitComparePageRoute(vGroups ...*gin.RouterGroup) {
	initComparePageV1Group(vGroups[0])
}

func initComparePageV1Group(v1Group *gin.RouterGroup) {
	compareController := v1.DefaultCompareController()
	compare := v1Group.Group(constants.Compare)
	{
		compare.GET("", compareController.GetCompareDetails)
		compare.GET(constants.PathSplitter+constants.List, compareController.GetCompareList)
	}
}
