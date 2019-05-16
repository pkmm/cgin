package controller

import (
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	UserId uint64
	respData map[string]interface{}
}

func (b *BaseController) GetAuthUserId(c *gin.Context) {
	val, ok := c.Get("uid")
	if !ok {
		service.SendResponse(c, errno.UserNotAuth, nil)
		return
	}

	userId, ok := val.(uint64)
	if !ok || userId == 0 {
		service.SendResponse(c, errno.UserNotAuth, nil)
		return
	}
	b.UserId = userId
}