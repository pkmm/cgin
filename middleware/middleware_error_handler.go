package middleware

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

func BusinessErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Abort()
				conf.AppLogger.Error("内部错误: => %#v", err)
				switch err.(type) {
				case *errno.BusinessErrorInfo:
					e := err.(*errno.BusinessErrorInfo)
					service.SendResponse(ctx, e, nil)
				case string, error:
					service.SendResponse(ctx, errno.InternalServerError, err)
				default:
					service.SendResponse(ctx, errno.InternalServerError, nil)
				}
			}
		}()
		ctx.Next()
	}
}
