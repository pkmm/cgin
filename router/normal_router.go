package router

import (
	"cgin/conf"
	v1 "cgin/controller/api/v1"
	"cgin/errno"
	"cgin/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func mapNormalRouter(rootRouter *gin.Engine) {
	// 通用
	rootRouter.Any("/", func(context *gin.Context) {
		currentAppEnv := conf.AppConfig.String(conf.AppEnvironment)
		service.SendResponse(context, errno.Welcome, fmt.Sprintf("current enviorment is [%s].", currentAppEnv))
	})

	// 未找到的路由
	rootRouter.NoRoute(func(context *gin.Context) {
		service.SendResponse(context, errno.NotSuchRouteException, nil)
	})

	// 普通的资源
	// api/
	apiNormal := rootRouter.Group(RootApiPrefix)
	{
		apiNormal.POST("/send_template_msg", v1.MiniProgramController.SendTemplateMsg)
		apiNormal.POST("/decode_verify_code", v1.VerifyCodeCtl.Recognize)
		apiNormal.GET("/daily/image", v1.DailyController.GetImage)
		apiNormal.GET("/daily/sentence", v1.DailyController.GetSentence)
	}
}
