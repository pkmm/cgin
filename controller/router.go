package controller

import (
	"github.com/gin-gonic/gin"
	"pkmm_gin/middleware"
	"pkmm_gin/util"
)

func MapRoute() *gin.Engine {

	ret := gin.New()
	ret.Use(middleware.MyRecovery(), middleware.ErrorHandle)

	ret.Static("/static", "/static")

	api := ret.Group(util.PathAPI)
	{
		api.Use(middleware.Auth)
		api.POST("/login", loginAction)
		api.POST("/get_scores", getScoresAction)
		api.POST("/set_account", setAccountAction)
	}

	return ret
}
