package controller

import (
	"cgin/conf"
	"cgin/middleware"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MapRoute() *gin.Engine {

	ret := gin.New()
	if "prod" == conf.AppConfig.DefaultString("appEnv", "prod") {
		ret.Use(middleware.MyRecovery(), middleware.ErrorHandle) // 正式环境不暴露出error
	} else {
		ret.Use(gin.Recovery()) // 开发的时候测试使用，可以比较方便的看到log
	}

	// 静态文件的目录
	ret.StaticFS("/static", http.Dir("static"))

	api := ret.Group(util.PathAPI)
	{
		api.Use(middleware.Auth)
		api.POST("/login", loginAction)
		api.POST("/get_scores", getScoresAction)
		api.POST("/set_account", setAccountAction)
		api.POST("/check_token", checkTokenAction)
	}

	return ret
}
