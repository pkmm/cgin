package system

import (
	"cgin/global"
	"cgin/model/common/resposne"
	"cgin/model/system"
	"cgin/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
	"go.uber.org/zap"
)

type SystemApi struct {
	// TODO：权限检查
}

func (s *SystemApi) Index(c *gin.Context) {
	name := c.Param("name")
	user := deliAutoSignService.GetUserByName(name)
	if user == nil {
		resposne.OkWithMsg("用户不存在", c)
		return
	}
	if err, html := deliAutoSignService.SignOne(user); err == nil {
		c.Data(200, "text/html", []byte(html))
		// 如果设置了 xwpusher的通知UID，发送微信通知
		if len(user.Uid) > 0 {
			global.WorkerPool.Submit(func() {
				msg := model.NewMessage(global.Config.Wxpusher.AppToken).
					SetSummary("签到成功: " + user.Username).
					SetContent(html).
					SetContentType(2).
					AddUId(user.Uid)
				wxpusher.SendMessage(msg)
			})
		}
		return
	} else {
		global.GLog.Error("未找到指定的用户", zap.Any("username:", name))
		resposne.FailWithMsg(err.Error(), c)
	}
}

func (s *SystemApi) WXPushCallBack(c *gin.Context) {
	var cbInfo request.WxpushCb
	if err := c.ShouldBindJSON(&cbInfo); err != nil {
		global.GLog.Error("wxpush回调参数验证失败", zap.Any("error", err))
		resposne.FailWithMsg("参数校验失败", c)
		return
	}
	global.GLog.Info("收到wxpush的回调信息", zap.Any("data", cbInfo))

	err := deliAutoSignService.UpdateUserWxpushUID(cbInfo.Data.Extra, cbInfo.Data.Uid)
	if err != nil {
		global.GLog.Error("更新用户的uid失败！", zap.Any("error", err))
	}
}

// GenerateQRCode 创建个人的二维码
func (s *SystemApi) GenerateQRCode(c *gin.Context) {
	name := c.Param("name")
	qrcode := model.Qrcode{AppToken: global.Config.Wxpusher.AppToken, Extra: name}
	resp, err := wxpusher.CreateQrcode(&qrcode)
	if err != nil {
		resposne.FailWithMsg("创建二维码失败！", c)
		return
	}
	resposne.OkWithData(resp, c)
}

func (s *SystemApi) DeliLogin(c *gin.Context) {
	reqData := system.DeliLoginData{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		resposne.FailWithMsg("解析参数失败!", c)
		return
	}
	if err, data := deliAutoSignService.Login(reqData.Mobile, reqData.Password); err == nil {
		// 登陆成功，更新token

		if deliAutoSignService.UserExists(reqData.Mobile) { // 用户已经存在的
			if err = deliAutoSignService.UpdateUserToken(reqData.Mobile, data.Data.Token); err == nil {
				resposne.OkWithMsg("操作成功！", c)
				return
			} else {
				resposne.FailWithMsg("更新用户token,操作失败！", c)
				return
			}
		} else { // 用户还不存在，是第一次登陆到系统，需要创建一个user
			if err = deliAutoSignService.CreateUser(reqData.Mobile, data.Data.Token, false); err == nil {
				resposne.OkWithMsg("创建用户成功！", c)
				return
			} else {
				resposne.FailWithMsg("创建用户失败！", c)
				return
			}
		}
	} else {
		resposne.FailWithMsg("登陆签到平台失败！" + err.Error(), c)
		return
	}
}

func (s *SystemApi) SetAutoSign(c *gin.Context) {
	name := c.Param("name")
	sign := c.GetBool("autoSign")
	if err := deliAutoSignService.SetAutoSign(name, sign); err != nil {
		resposne.FailWithMsg("操作失败！", c)
		return
	}
	resposne.OkWithData("操作成功！", c)
}