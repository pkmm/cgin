package middleware

import (
	"cgin/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GinRecovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				context.AbortWithStatus(http.StatusInternalServerError)
				global.GLog.Error("server panic: ", zap.Any("Panic", err))
			}
		}()
		context.Next()
	}
}
