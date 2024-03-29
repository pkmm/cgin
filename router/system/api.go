package system

import (
	v1 "cgin/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (a *ApiRouter) InitApiRouter(r *gin.RouterGroup) {
	apiRouter := r.Group("user")
	var api = v1.ApiGroupApp.SystemApiGroup.SystemApi
	{
		apiRouter.GET("sign/:name", api.Index)
		apiRouter.POST("wxpusher/cb", api.WXPushCallBack)
		apiRouter.GET("qrcode/:name", api.GenerateQRCode)
		apiRouter.POST("deliLogin", api.DeliLogin)
		apiRouter.Any("setAutoSign/:name", api.SetAutoSign)
	}
}
