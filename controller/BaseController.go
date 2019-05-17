package controller

import (
	"cgin/errno"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BaseController struct {
	UserId uint64
	Params map[string]interface{}
}

func (b *BaseController) GetAuthUserId(c *gin.Context) {
	val, ok := c.Get("uid")
	if !ok {
		panic(errno.UserNotAuth)
	}

	userId, ok := val.(uint64)
	if !ok || userId == 0 {
		panic(errno.UserNotAuth)
	}
	b.UserId = userId
}

// 请求的中json参数解析到params
func (b *BaseController) Init(c *gin.Context) {
	b.Params = map[string]interface{}{}
	if err := c.ShouldBindWith(&b.Params, binding.JSON); err != nil {
		panic(errno.InvalidParameters.AppendErrorMsg(err.Error()))
	}
}