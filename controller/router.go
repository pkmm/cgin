package controller

import (
	"github.com/gin-gonic/gin"
	"pkmm_gin/middleware"
	"pkmm_gin/util"
)

func MapRoute() *gin.Engine {

	ret := gin.New()
	ret.Use(gin.Recovery())

	ret.Static("/static", "./static")

	api := ret.Group(util.PathAPI)
	api.Use(middleware.Auth)
	api.POST("/login", loginAction)

	return ret
}
