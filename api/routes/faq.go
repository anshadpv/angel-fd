package routes

import (
	v1 "github.com/angel-one/fd-core/api/v1"
	"github.com/angel-one/fd-core/constants"
	"github.com/gin-gonic/gin"
)

func initFAQ(vGroups ...*gin.RouterGroup) {
	initFAQV1Group(vGroups[0])
}

func initFAQV1Group(v1Group *gin.RouterGroup) {
	faqController := v1.DefaultFAQController()
	faqs := v1Group.Group(constants.FAQ)
	{
		faqs.GET(constants.PathParam+constants.Tag, faqController.GetFAQs)
	}
}
