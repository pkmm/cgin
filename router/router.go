package router

import (
	"cgin/conf"
	"cgin/controller"
	"cgin/errno"
	"cgin/middleware"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MapRoute() *gin.Engine {

	router := gin.New()
	// 全局的中间件
	if "prod" == conf.AppConfig.DefaultString("appEnv", "prod") {
		router.Use(middleware.MyRecovery(), middleware.RequestLogger) // 正式环境不暴露出error
	} else {
		router.Use(gin.Recovery()) // 开发的时候测试使用，可以比较方便的看到log
	}

	// 通用
	router.Any("/", func(context *gin.Context) {
		currentAppEnv := conf.AppConfig.String("appEnv")
		service.SendResponse(context, errno.Welcome, currentAppEnv)
	})

	// 未找到的路由
	router.NoRoute(func(context *gin.Context) {
		service.SendResponse(context, errno.NotSuchRouteException, nil)
	})

	// 静态文件的目录
	router.StaticFS("/static", http.Dir("static"))

	// 业务的逻辑 API
	apiAuth := router.Group(util.PathAPI).Use(middleware.Auth)
	{
		// 需要进行认证的业务API
		apiAuth.POST("/get_scores", controller.UserController.GetScoresAction)
		apiAuth.POST("/set_account", controller.UserController.SetAccountAction)
		apiAuth.POST("/check_token", controller.UserController.CheckTokenAction)
		apiAuth.POST("/send_template_msg", controller.UserController.SendTemplateMsg)
	}

	apiNotAuth := router.Group(util.PathAPI)
	{
		// 不需要进行认证的API
		apiNotAuth.POST("/login", controller.UserController.LoginAction)
		apiNotAuth.POST("/decode_verify_code", controller.VerifyCodeCtl.Recognize)
	}

	return router
}
