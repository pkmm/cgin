package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapAuthRouter(router *gin.Engine) {
	apiAuth := router.Group(AuthPrefix).Use(middleware.Auth)
	{
		apiAuth.POST("/me", v1.AuthController.Me)
	}

	apiNotAuth := router.Group(AuthPrefix)
	{
		apiNotAuth.POST("/login", v1.AuthController.Login)
	}
}
