package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func initScoresRouter(scoreRouter *gin.Engine) {
	api := scoreRouter.Group(Score).Use(middleware.Auth)
	{
		api.GET("/", v1.Student.GetScores)
	}
}
