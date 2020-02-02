package middleware

import (
	"cgin/conf"
	"github.com/gin-gonic/gin"
)

func RequestLogger(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		conf.AppLogger.Info("Request: URL[%s], RemoteIP[%s]", cCp.Request.URL.Path, cCp.Request.RemoteAddr)
	}()
	c.Next()
	if len(c.Errors) != 0 {
		go func() {
			conf.AppLogger.Error("Server errors: " + cCp.Errors.String())
		}()
	}
}
