package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitTokenRoute(vGroups ...*gin.RouterGroup) {
	initTokenV1Group(vGroups[0])
}

func initTokenV1Group(v1Group *gin.RouterGroup) {
	tokenController := v1.DefaultTokenController()

	token := v1Group.Group(constants.Token)
	{
		token.GET(constants.PathParam+constants.Provider, tokenController.GetToken)
	}
}
