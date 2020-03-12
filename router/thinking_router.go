package router

import (
	"cgin/controller"
	"github.com/gin-gonic/gin"
)

func mapThinkingRouter(router *gin.Engine) {
	thinkingApi := router.Group(Thinking)
	{
		thinkingApi.Any("/list", controller.ThinkingController.GetList)
	}
}
