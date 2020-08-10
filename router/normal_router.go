package router

import (
	"cgin/conf"
	v1 "cgin/controller/api/v1"
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

func initNormalRouter(normalRouter *gin.Engine) {
	// 通用
	normalRouter.Any("/", func(context *gin.Context) {
		currentAppEnv := conf.AppEnvironment()
		info := struct {
			Github string `json:"github"`
			Env    string `json:"env"`
			Author string `json:"author"`
		}{
			Github: "https://github.com/pkmm",
			Env:    currentAppEnv,
			Author: "相顾无言",
		}
		service.SendResponse(
			context,
			errno.Welcome,
			info,
		)
	})

	// 未找到的路由
	normalRouter.NoRoute(func(context *gin.Context) {
		service.SendResponse(context, errno.NotSuchRouteException, nil)
	})

	apiNormal := normalRouter.Group(RootApiPrefix)
	{
		apiNormal.POST("/decode_verify_code", v1.VerifyCodeCtl.Recognize)
		apiNormal.GET("/daily/images", v1.DailyController.GetImage)
		apiNormal.GET("/daily/sentences", v1.DailyController.GetSentence)
	}
}
