package router

import (
	"cgin/controller/api/v1"
	"github.com/gin-gonic/gin"
)

func mapThinkingRouter(router *gin.Engine) {
	thinkingApi := router.Group(Thinking)
	{
		thinkingApi.Any("/list", v1.ThinkingController.GetList)
	}
}
