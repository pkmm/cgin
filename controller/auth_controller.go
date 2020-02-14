package controller

import (
	"cgin/conf"
	"cgin/constant/devicetype"
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type authController struct {}

var AuthController = &authController{}

// @Summary 登录
// @Produce json
// @Accept json
// @Router /auth/login [post]
// @Success 200 {object} service.Response
// @Param auth_model body co.AuthModel true "auth model"
func (a *authController) Login(c *gin.Context) {
	var (
		openid, sign, token string
		deviceType          devicetype.DeviceType
		err                 error
	)
	helper := context_helper.New(c)
	deviceType = devicetype.DeviceType(helper.GetInt("device_type"))
	openid = helper.GetString("openid")
	openid = helper.GetString("sign")

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
		helper.Response(data)
	}
}

// @Summary 获取认证的自己
// @Accept json
// @Produce json
// @Router /auth/me [post]
// @Param auth_credit body co.AuthCredit true "get auth self"
// @Success 200 {object} service.Response
func (a *authController) Me(c *gin.Context) {
	helper := context_helper.New(c)
	user := service.User.GetUser(helper.GetAuthUserId())
	data := gin.H{
		"user": user,
	}
	helper.Response(data)
}
