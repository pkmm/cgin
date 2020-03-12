package router

import (
	"cgin/controller"
	"github.com/gin-gonic/gin"
)

func mapScoresRouter(router *gin.Engine) {
	api := router.Group(Score)
	{
		api.GET("/*studentId", controller.Student.GetScores)
	}
}
