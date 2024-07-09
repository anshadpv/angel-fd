package routes

import (
	"github.com/angel-one/fd-core/api/events"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func InitWebhookRoutes(vGroups ...*gin.RouterGroup) {
	initWebhookRoutes(vGroups[0])
}

func initWebhookRoutes(group *gin.RouterGroup) {
	controller := events.DefaultWebhooksController()
	group.POST(constants.PathSplitter+constants.UpSwingProvider+constants.PathSplitter+constants.UpSwingWebhookPath, controller.ReadUpSwingMessage)
}
