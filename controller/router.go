package controller

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/middleware"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MapRoute() *gin.Engine {

	router := gin.New()
	if "prod" == conf.AppConfig.DefaultString("appEnv", "prod") {
		router.Use(middleware.MyRecovery(), middleware.RequestLogger) // 正式环境不暴露出error
	} else {
		router.Use(gin.Recovery()) // 开发的时候测试使用，可以比较方便的看到log
	}

	// 通用
	router.Any("/", func(context *gin.Context) {
		service.SendResponse(context, errno.Welcome, nil)
	})

	// 未找到的路由
	router.NoRoute(func(context *gin.Context) {
		service.SendResponse(context, errno.NotSuchRouteException, nil)
	})

	// 静态文件的目录
	router.StaticFS("/static", http.Dir("static"))

	api := router.Group(util.PathAPI)
	{
		api.Use(middleware.Auth)
		api.POST("/login", loginAction)
		api.POST("/get_scores", getScoresAction)
		api.POST("/set_account", setAccountAction)
		api.POST("/check_token", checkTokenAction)
		api.POST("/send_template_msg", sendTemplateMsg)
	}

	return router
}
