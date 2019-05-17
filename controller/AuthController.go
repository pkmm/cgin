package controller

import (
	"cgin/conf"
	"cgin/errno"
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
	var (
		openid, sign, token string
		ok                  bool
		err                 error
	)
	arg := map[string]interface{}{}
	if err = c.ShouldBindWith(&arg, binding.JSON); err != nil {
		panic(errno.InvalidParameters)
	}

	if openid, ok = arg["openid"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数openid必须提供"))
	}
	if sign, ok = arg["sign"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数sign必须提供"))
	}

	sign2 := util.Md5String("xiaocc_ai_liu_yan_lin" + conf.AppConfig.String("miniprogram_app_id") + openid)
	if sign != sign2 {
		panic(errno.InvalidParameters.AppendErrorMsg("签名验证失败"))
	}

	user := service.AuthService.LoginFromMiniProgram(openid)
	if token, err = service.JWTSrv.GenerateToken(user); err != nil {
		panic(errno.GenerateJwtTokenFailed.AppendErrorMsg(err.Error()))
	}
	data := gin.H{
		"user":  user,
		"token": token,
	}
	service.SendResponseSuccess(c, data)
}

func (a *authController) Me(c *gin.Context) {
	a.BaseController.GetAuthUserId(c)
	user := service.User.GetUser(a.UserId)
	data := gin.H{
		"user": user,
	}
	service.SendResponseSuccess(c, data)
}
