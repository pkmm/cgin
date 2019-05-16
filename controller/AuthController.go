package controller

import (
	"cgin/conf"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type authController struct {
	BaseController
}

var AuthController = &authController{}

func (a *authController) Login(c *gin.Context) {
	arg := map[string]interface{}{}
	if err := c.ShouldBindWith(&arg, binding.JSON); err != nil {
		service.SendResponseWithInvalidParameters(c, "参数不合法")
		return
	}
	var (
		openid, sign string
		ok           bool
	)
	if openid, ok = arg["openid"].(string); !ok {
		service.SendResponseWithInvalidParameters(c, "参数openid必须提供")
		return
	}
	if sign, ok = arg["sign"].(string); !ok {
		service.SendResponseWithInvalidParameters(c, "参数sign必须提供")
		return
	}

	sign2 := util.Md5String("xiaocc_ai_liu_yan_lin" + conf.AppConfig.String("miniprogram_app_id") + openid)
	if sign != sign2 {
		service.SendResponseWithInvalidParameters(c, "参数验证失败sign值错误")
		return
	}

	user := service.AuthService.LoginFromMiniProgram(openid)
	data := map[string]interface{}{
		"user": user,
	}
	service.SendResponseSuccess(c, data)
}

func (a *authController) Me(c *gin.Context) {
	a.BaseController.GetAuthUserId(c)
	user := service.User.GetUser(a.UserId)
	a.respData["user"] = user
	service.SendResponseSuccess(c, data)
}
