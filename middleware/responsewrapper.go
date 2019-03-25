package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pkmm_gin/conf"
	"pkmm_gin/errno"
	"pkmm_gin/service"
)

func ErrorHandle(c *gin.Context) {
	c.Next()

	if len(c.Errors) != 0 {
		fmt.Println(c.Errors)
		service.SendResponse(c, errno.InternalServerError, nil)
	}
}

func MyRecovery() gin.HandlerFunc {

	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 发生panic 在这里handle处理
				conf.AppLogger.Error("panic, ", err)
				service.SendResponse(context, errno.InternalServerError, nil)
				context.Abort()
			}
		}()
		context.Next()
	}
}
