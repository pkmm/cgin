package router

import (
	"cgin/conf"
	"cgin/controller"
	"cgin/errno"
	"cgin/middleware"
	"cgin/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RootApiPrefix     = "/api"
	AuthPrefix        = "/api/auth"
	StudentPrefix     = "/api/student"
	MiniProgramPrefix = RootApiPrefix + "/mini_program"
	Trigger = "/api/trigger/"
)

func MapRoute() *gin.Engine {

	router := gin.New()
	// 全局的中间件
	if "prod" == conf.AppConfig.DefaultString("appEnv", "prod") {
		router.Use(gzip.Gzip(gzip.DefaultCompression), middleware.BusinessErrorHandler(), middleware.RequestLogger) // 正式环境不暴露出error
	} else {
		router.Use(gzip.Gzip(gzip.DefaultCompression), middleware.BusinessErrorHandler()) // 开发的时候测试使用，可以比较方便的看到log
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
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

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
		apiNormal.POST("/send_template_msg", controller.MiniProgramController.SendTemplateMsg)
		apiNormal.POST("/decode_verify_code", controller.VerifyCodeCtl.Recognize)
		apiNormal.GET("/daily_image", controller.DailyController.GetImage)
		apiNormal.GET("/daily_sentence", controller.DailyController.GetSentence)
	}

	// 小程序
	// api/mini_program
	apiMiniProgram := router.Group(MiniProgramPrefix)
	{
		// 不需要认证的
		apiMiniProgram.POST("/get_index_preference", controller.MiniProgramController.GetIndexPreference)
		apiMiniProgram.POST("/set_index_config", controller.MiniProgramController.SetIndexConfig)
		apiMiniProgram.POST("/get_notifications", controller.MiniProgramController.GetNotification)
		apiMiniProgram.POST("/get_sponsors", controller.MiniProgramController.GetSponsors)

		// 以下的API 需要认证
		apiMiniProgram.POST("/config_menu", controller.MiniProgramController.DisposeMenu).Use(middleware.Auth)
		apiMiniProgram.POST("/change_notification", controller.MiniProgramController.UpdateOrCreateNotification)
		apiMiniProgram.POST("/get_hermann_memorial", controller.HermannRememberController.GetTodayTaskInfo)
		apiMiniProgram.POST("/add_hermann_memorial", controller.HermannRememberController.SaveUserRememberTask)
	}

	apiTrigger := router.Group(Trigger)
	{
		if conf.AppConfig.String("appEnv") != "dev" {
			apiTrigger.Use(middleware.Auth)
		}
		apiTrigger.Any("/cron", controller.CronTaskController.TriggerTask)
	}

	return router
}
