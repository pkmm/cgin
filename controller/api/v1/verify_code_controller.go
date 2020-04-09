package v1

import (
	"cgin/controller/contextHelper"
	"cgin/errno"
	"cgin/zcmu"
	"github.com/gin-gonic/gin"
	"image/gif"
)

type verifyCodeController struct{}

var VerifyCodeCtl = &verifyCodeController{}

// @Summary 验证码识别
// @Tags 验证码识别
// @Param img formData file true "image of verify code"
// @Router /decode_verify_code [post]
// @Success 200 object service.Response
// @Produce json
// @Accept image/gif
func (v *verifyCodeController) Recognize(c *gin.Context) {
	helper := contextHelper.New(c)
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
	helper.Response(gin.H{
		"text": text,
	})
}
