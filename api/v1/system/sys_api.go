package system

import (
	"cgin/global"
	"cgin/model/common/resposne"
	"cgin/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApi struct {
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

	//deliAutoSignService.UpdateUserWxpushUID()

}
