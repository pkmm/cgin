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
				switch err.(type) {
				case *errno.BusinessErrorInfo:
					e := err.(*errno.BusinessErrorInfo)
					service.SendResponse(ctx, e, nil)
				case string, error:
					conf.Logger.Error("服务器错误:[%s] => %#v", ctx.Request.RequestURI, err)
					service.SendResponse(ctx, errno.InternalServerError, err)
				default:
					conf.Logger.Error("服务器错误:[%s] => %#v", ctx.Request.RequestURI, err)
					service.SendResponse(ctx, errno.InternalServerError, nil)
				}
			}
		}()
		ctx.Next()
	}
}
