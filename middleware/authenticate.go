package middleware

import (
	"bytes"
	"cgin/conf"
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
)

// 使用JWT进行认证

type authData struct {
	Token string `json:"token"`
}

func Auth(c *gin.Context) {
	// 创建buffer 备份body
	buf, _ := ioutil.ReadAll(c.Request.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

	data := &authData{}
	c.Request.Body = rdr1
	if err := c.ShouldBindWith(data, binding.JSON); err != nil {
		conf.AppLogger.Error("login request failed. " + err.Error())
		c.Request.Body = rdr2
		service.SendResponse(c, errno.InvalidParameters, nil)
		c.Abort()
		return
	}
	c.Request.Body = rdr2

	if data.Token == "" {
		service.SendResponse(c, errno.TokenNotValid, nil)
		c.Abort()
		return
	}

	claims, err := service.JWTSrv.GetAuthClaims(data.Token)
	if err != nil {
		service.SendResponse(c, errno.TokenNotValid, nil)
		c.Abort()
		return
	}

	c.Set("uid", claims.Uid)
	c.Next()

	return
}
