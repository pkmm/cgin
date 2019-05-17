package router

import (
	"cgin/conf"
	"cgin/controller"
	"cgin/errno"
	"cgin/middleware"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)


const (
	RootApiPrefix = "/api"
	AuthPrefix = "/api/auth"
	StudentPrefix = "/api/student"
)

func MapRoute() *gin.Engine {

	router := gin.New()
	// 全局的中间件
	if "prod" == conf.AppConfig.DefaultString("appEnv", "prod") {
		router.Use(middleware.BusinessErrorHandler(), middleware.RequestLogger) // 正式环境不暴露出error
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

	// 认证业务的逻辑 API
	// api/auth
	apiAuth := router.Group(AuthPrefix).Use(middleware.Auth)
	{
		// 需要进行认证的业务API
		apiAuth.POST("/me", controller.AuthController.Me)
	}

	// api/auth
	apiNotAuth := router.Group(AuthPrefix)
	{
		// 不需要进行认证的API
		apiNotAuth.POST("/login", controller.AuthController.Login)
	}

	// api/student
	apiStudent := router.Group(StudentPrefix).Use(middleware.Auth)
	{
		apiStudent.POST("/", controller.Student.GetStudent)
		apiStudent.POST("/scores", controller.Student.GetScores)
		apiStudent.POST("/update_edu_account", controller.Student.UpdateEduAccount)
	}

	// 普通的资源
	// api/
	apiNormal := router.Group(RootApiPrefix)
	{
		apiNormal.POST("/send_template_msg", controller.UserController.SendTemplateMsg)
		apiNormal.POST("/decode_verify_code", controller.VerifyCodeCtl.Recognize)
	}

	return router
}
