package controller

import (
	"cgin/errno"
	"cgin/service"
	"cgin/zcmu"
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
		panic(errno.InvalidParameters.AppendErrorMsg(err.Error()))
	}
	openedFile, err := file.Open()
	if err != nil {
		panic(errno.InvalidParameters.AppendErrorMsg(err.Error()))
	}
	img, err := gif.Decode(openedFile)
	if err != nil {
		panic(errno.InvalidParameters.AppendErrorMsg(err.Error()))
	}
	text, err = zcmu.Predict(img)
	if err != nil {
		panic(errno.InvalidParameters.AppendErrorMsg(err.Error()))
	}
	service.SendResponse(c, errno.Success, map[string]string{
		"text": text,
	})
}
