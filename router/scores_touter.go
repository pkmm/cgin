package router

import (
	"cgin/controller"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapScoresRouter(router *gin.Engine) {
	api := router.Group(Score).Use(middleware.Auth)
	{
		api.GET("/", controller.Student.GetScores)
	}
}
