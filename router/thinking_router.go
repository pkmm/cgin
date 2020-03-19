package router

import (
	"cgin/controller/api/v1"
	"github.com/gin-gonic/gin"
)

func initThinkingRouter(thinkRouter *gin.Engine) {
	thinkingApi := thinkRouter.Group(Thinking)
	{
		thinkingApi.Any("/list", v1.ThinkingController.GetList)
	}
}
