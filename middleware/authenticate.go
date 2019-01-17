package middleware

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"pkmm_gin/errno"
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
	// 创建buffer 备份body
	buf, _ := ioutil.ReadAll(c.Request.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
	data := &authData{}
	c.Request.Body = rdr1
	if err := c.BindJSON(data); err != nil {
		logs.Error("login request failed. " + err.Error())
		c.Request.Body = rdr2
		c.Abort()
		return
	}

	user := service.User.CheckAndGetUserByUserIdAndAccessToken(data.UserId, data.AccessToken)
	if user == nil {
		c.Request.Body = rdr2
		service.SendResponse(c, errno.ErrUserNotFound, nil)
		c.Abort()
		return
	}
	c.Set("user", user)
	c.Request.Body = rdr2

	c.Next()

	return
}
