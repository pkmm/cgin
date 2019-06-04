package controller

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type weChatController struct {
	BaseController
}

var WeChatController = &weChatController{}

func (w *weChatController) SendTemplateMsg(c *gin.Context) {
	w.ProcessParams(c)

	formId, ok := w.Params["form_id"].(string)
	if !ok {
		panic(errno.NormalException.ReplaceErrorMsgWith("require form id."))
	}
	openId, ok := w.Params["open_id"].(string)
	if !ok {
		panic(errno.NormalException.ReplaceErrorMsgWith("require open id."))
	}
	templateKeyData := &util.TemplateMsgData{}
	templateKeyData.Keyword1.Value = "11"
	templateKeyData.Keyword2.Value = "22"
	msg := &util.TemplateMsg{
		FormId:     formId,
		ToUser:     openId,
		TemplateId: conf.AppConfig.String("template_id"),
		Page:       conf.AppConfig.String("template_msg_open_page"),
		Data:       templateKeyData,
	}
	ret := util.SendUserTemplateMsg(msg)
	w.Response(c, ret)
}
