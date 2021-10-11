package system

import (
	"cgin/global"
	"cgin/model/common/resposne"
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
