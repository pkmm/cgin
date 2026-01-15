package middleware

import (
	"cgin/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
