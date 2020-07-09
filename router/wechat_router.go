package router

import (
	v1 "cgin/controller/api/v1"
	"github.com/gin-gonic/gin"
)

func initWeChatRouter(r *gin.Engine) {
	wx := r.Group(Wx)
	{
		wx.Any("", v1.WeChatController.Index)
	}
}
