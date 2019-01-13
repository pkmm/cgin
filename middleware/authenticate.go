package middleware

import (
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"pkmm_gin/service"
	"strings"
)

type authData struct {
	UserId      uint64 `json:"__user_id"`
	AccessToken string `json:"__access_token"`
}

func Auth(c *gin.Context) {
	if strings.Index(c.Request.URL.Path, "login") != -1 {
		c.Next()
		return
	}
	data := &authData{}
	if err := c.Bind(data); err != nil {
		logs.Error("login request failed. " + err.Error())
		c.Abort()
		return
	}

	user := service.User.CheckAndGetUserByUserIdAndAccessToken(data.UserId, data.AccessToken)
	if user == nil {
		c.Abort()
		return
	}
	c.Set("user", &user)
	c.Next()

	return
}
