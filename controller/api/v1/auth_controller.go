package v1

import (
	"cgin/conf"
	"cgin/constant/devicetype"
	"cgin/controller/contextHelper"
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"cgin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type authController struct{}

var AuthController = &authController{}

// @Summary 登录
// @Produce json
// @Accept json
// @Router /auth/login [post]
// @Success 200 object service.Response
// @Param auth_model body co.AuthModel true "认证的信息"
func (a *authController) Login(c *gin.Context) {
	var (
		openid, sign, token string
		loginDeviceType     devicetype.DeviceType
		err                 error
		user                *model.User
	)

	helper := contextHelper.New(c)
	loginDeviceType = devicetype.DeviceType(helper.GetInt("device_type"))

	switch loginDeviceType {
	case devicetype.MiniProgram:
		openid = helper.GetString("openid")
		sign = strings.ToUpper(helper.GetString("sign"))
		sign2 := util.Md5String(conf.AppConfig.String("normal.random.str") + conf.AppConfig.String("miniprogram_app_id") + openid)
		if sign != sign2 {
			if conf.IsDev() {
				fmt.Println("签名：", sign2)
			}
			panic(errno.InvalidParameters.AppendErrorMsg("签名验证失败"))
		}
		err, user = service.AuthService.LoginFromMiniProgram(openid)
	case devicetype.WebBrowser:
		err, user = service.AuthService.LoginFromWebBrowser(
			helper.GetString("username"), helper.GetString("password"))
	default:
		panic(errno.NormalException.ReplaceErrorMsgWith("未支持的device type"))
	}

	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith("登录失败:" + err.Error()))
	}

	if token, err = service.JWTSrv.GenerateToken(user); err != nil {
		panic(errno.GenerateJwtTokenFailed.AppendErrorMsg(err.Error()))
	}

	helper.Response(gin.H{
		"user":  user,
		"token": token,
	})
}

// @Summary 获取认证的自己
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Router /auth/me [post]
// @Success 200 object service.Response
func (a *authController) Me(c *gin.Context) {
	helper := contextHelper.New(c)
	err, user := model.GetUserById(helper.GetAuthUserId())
	if err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	helper.Response(gin.H{
		"user": user,
	})
}
