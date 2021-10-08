package initialize

import (
	"cgin/global"
	"cgin/router"
	"github.com/gin-gonic/gin"
	"net/http"
)
// 初始化所有的路由
func Routers() *gin.Engine {
	var Router = gin.Default()

	// 静态文件的地址
	Router.StaticFS(global.G_CONFIG.Local.Path, http.Dir(global.G_CONFIG.Local.Path))

	// 获取响应的路由实例
	systemRouter := router.RouterGroupApp.System
	PublicRouter := Router.Group("")
	{
		// 健康检测
		PublicRouter.GET("/health", func(context *gin.Context) {
			context.JSON(200, "ok")
		})
	}
	PrivateRouter := Router.Group("")

	{
		systemRouter.InitApiRouter(PrivateRouter)
	}
	global.G_LOG.Info("router register success")
	return Router
}
