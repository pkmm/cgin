package middleware

import (
	"cgin/conf"
	"github.com/gin-gonic/gin"
)

func RequestLogger(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		conf.Logger.Info("Request: URL[%s], RemoteIP[%s]", cCp.Request.RequestURI, cCp.ClientIP())
	}()
	c.Next()
	if len(c.Errors) != 0 {
		go func() {
			conf.Logger.Error("Server errors: " + cCp.Errors.String())
		}()
	}
}
