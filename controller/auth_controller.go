package controller

import (
	"cgin/conf"
	"cgin/constant/devicetype"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type authController struct {
	BaseController
}

var AuthController = &authController{}

func (a *authController) Login(c *gin.Context) {
	var (
		openid, sign, token string
		deviceType          devicetype.DeviceType
		ok                  bool
		err                 error
	)
	a.processParams(c)
	deviceType = devicetype.DeviceType(a.Params["device_type"].(float64))
	if openid, ok = a.Params["openid"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数openid必须提供"))
	}
	if sign, ok = a.Params["sign"].(string); !ok {
		panic(errno.InvalidParameters.AppendErrorMsg("参数sign必须提供"))
	}

	sign2 := util.Md5String(conf.AppConfig.String("normal.random.str") + conf.AppConfig.String("miniprogram_app_id") + openid)
	if sign != sign2 {
		panic(errno.InvalidParameters.AppendErrorMsg("签名验证失败"))
	}
	if deviceType == devicetype.MiniProgram {
		user := service.AuthService.LoginFromMiniProgram(openid)
		if token, err = service.JWTSrv.GenerateToken(user); err != nil {
			panic(errno.GenerateJwtTokenFailed.AppendErrorMsg(err.Error()))
		}
		data := gin.H{
			"user":  user,
			"token": token,
		}
		a.response(c, data)
	}

}

func (a *authController) Me(c *gin.Context) {
	a.getAuthUserId(c)
	user := service.User.GetUser(a.UserId)
	data := gin.H{
		"user": user,
	}
	a.response(c, data)
}