package middleware

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BusinessErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Abort()
				switch err.(type) {
				case *errno.BusinessErrorInfo:
					e := err.(*errno.BusinessErrorInfo)
					switch e.Code {
					// 设置状态码
					case errno.TokenNotValid.Code:
						service.SendResponseWithStatus(ctx, e, nil, http.StatusUnauthorized)
					case errno.PermissionDenied.Code:
						service.SendResponseWithStatus(ctx, e, nil, http.StatusForbidden)
					default:
						service.SendResponse(ctx, e, nil)
					}
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
