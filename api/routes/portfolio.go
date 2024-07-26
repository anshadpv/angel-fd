package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitPortfolioRoute(vGroups ...*gin.RouterGroup) {
	initPortfolioV1Group(vGroups[0])
}

func initPortfolioV1Group(v1Group *gin.RouterGroup) {
	portfolioController := v1.DefaultPortfolioController()

	portfolio := v1Group.Group(constants.Portfolio)
	{
		portfolio.GET("", portfolioController.GetPortfolio)
		portfolio.GET(constants.PathParam+constants.Provider+constants.PathSplitter+constants.Networth, portfolioController.GetNetworth)
		portfolio.GET(constants.Details, portfolioController.GetPortfolioDetails)
	}
}
