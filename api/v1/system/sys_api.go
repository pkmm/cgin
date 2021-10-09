package system

import (
	"github.com/gin-gonic/gin"
)

type SystemApi struct {
}

func (s *SystemApi) Index(c *gin.Context) {
	name := c.Param("name")
	user := deliAutoSignService.GetUserByName(name)
	if user == nil {
		c.JSON(200, "用户不存在")
		return
	}
	if err, html := deliAutoSignService.SignOne(user); err == nil {
		c.Data(200, "text/html", []byte(html))
		return
	} else {
		c.JSON(200, err)
	}
}
