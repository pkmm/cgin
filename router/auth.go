package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func initAuthRouter(authRouter *gin.Engine) {
	apiAuth := authRouter.Group(AuthPrefix).Use(middleware.Auth())
	{
		apiAuth.POST("/me", v1.AuthController.Me)
	}

	apiNotAuth := authRouter.Group(AuthPrefix)
	{
		apiNotAuth.POST("/login", v1.AuthController.Login)
	}
}
