package router

import (
	"cgin/conf"
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
}
