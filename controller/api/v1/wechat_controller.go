package v1

import (
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type wechatCtl struct{}

var WeChatController = &wechatCtl{}

func (w *wechatCtl) Index(context *gin.Context) {
	service.WeChatAppService.Serve(context.Writer, context.Request)
}
