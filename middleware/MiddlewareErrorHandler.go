package middleware

import (
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BusinessErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case *errno.BusinessErrorInfo:
					e := err.(*errno.BusinessErrorInfo)
					ctx.Abort()
					ctx.JSON(http.StatusOK, gin.H{
						"code": e.Code,
						"msg":  e.Msg,
						"data": nil,
					})
				default:
					//panic(err)
					service.SendResponse(ctx, errno.InternalServerError, nil)
				}
			}
		}()
		ctx.Next()
	}
}
