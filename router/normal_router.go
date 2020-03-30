package router

import (
	"cgin/conf"
	v1 "cgin/controller/api/v1"
	"cgin/errno"
	"cgin/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func initNormalRouter(normalRouter *gin.Engine) {
	// 通用
	normalRouter.Any("/", func(context *gin.Context) {
		currentAppEnv := conf.AppEnvironment()
		service.SendResponse(context, errno.Welcome, fmt.Sprintf("current enviorment is [%s].", currentAppEnv))
	})

	// 未找到的路由
	normalRouter.NoRoute(func(context *gin.Context) {
		service.SendResponse(context, errno.NotSuchRouteException, nil)
	})

	// 普通的资源
	// api/
	apiNormal := normalRouter.Group(RootApiPrefix)
	{
		apiNormal.POST("/decode_verify_code", v1.VerifyCodeCtl.Recognize)
		apiNormal.GET("/daily/images", v1.DailyController.GetImage)
		apiNormal.GET("/daily/sentences", v1.DailyController.GetSentence)
	}
}
