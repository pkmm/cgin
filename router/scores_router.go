package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapScoresRouter(router *gin.Engine) {
	api := router.Group(Score).Use(middleware.Auth)
	{
		api.GET("/", v1.Student.GetScores)
	}
}
