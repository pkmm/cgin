package middleware

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

func RequestLogger(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		conf.AppLogger.Info("Request: URL[%s], RemoteIP[%s]", cCp.Request.URL, cCp.Request.RemoteAddr)
	}()
	c.Next()
	if len(c.Errors) != 0 {
		go func() {
			conf.AppLogger.Error("Server errors: " + cCp.Errors.String())
		}()
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
