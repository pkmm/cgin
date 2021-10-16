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
	Router.StaticFS(global.Config.Local.Path, http.Dir(global.Config.Local.Path))

	// 默认页面的地址
	Router.Any("/img", func(c *gin.Context) {
		c.File("static/images/2019-06-10.webp")
	})

	Router.Any("", func(c *gin.Context) {
		c.String(200, `欢迎访问cgin服务！
得力签到：https://api.qwesde.com/user/sign/<username>
获取微信通知：https://api.qwesde.com/user/qrcode/<username>
设置是否使用自动签到：https://api.qwesde.com/user/setAutoSign/<username>?autoSign=<false|true>
登陆签到系统：https://api.qwesde.com/user/deliLogin [POST] {"mobile"": "<phone>", "password": "<password>"}
上述的<username>请更换成自己的名字
`)
	})

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
	global.GLog.Info("router register success")
	return Router
}
