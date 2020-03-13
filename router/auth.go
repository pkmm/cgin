package router

import (
	"cgin/controller"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapAuthRouter(router *gin.Engine) {
	apiAuth := router.Group(AuthPrefix).Use(middleware.Auth)
	{
		apiAuth.POST("/me", controller.AuthController.Me)
	}

	apiNotAuth := router.Group(AuthPrefix)
	{
		apiNotAuth.POST("/login", controller.AuthController.Login)
	}
}
