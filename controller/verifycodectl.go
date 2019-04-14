package controller

import (
	"cgin/errno"
	"cgin/service"
	"cgin/zcmuES"
	"github.com/gin-gonic/gin"
	"image/gif"
)

type verifyCodeController struct{}

var VerifyCodeCtl = &verifyCodeController{}

// 识别验证码中的文字
// 提供服务给外面调用
func (v *verifyCodeController) Recognize(c *gin.Context) {
	var text string
	file, err := c.FormFile("img")
	if err != nil {
		service.SendResponse(c, errno.InvalidParameters, err.Error())
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		service.SendResponse(c, errno.InvalidParameters, err.Error())
		return
	}
	img, err := gif.Decode(openedFile)
	if err != nil {
		service.SendResponse(c, errno.InvalidParameters, err.Error())
		return
	}
	text, err = zcmuES.Predict(img)
	if err != nil {
		service.SendResponse(c, errno.InvalidParameters, err.Error())
		return
	}
	service.SendResponse(c, errno.Success, map[string]string{
		"text": text,
	})
}
