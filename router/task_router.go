package router

import (
	"cgin/conf"
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapTaskRouter(router *gin.Engine) {
	apiTrigger := router.Group(Trigger)
	{
		if conf.AppConfig.String(conf.AppEnvironment) != conf.AppEnvDev {
			apiTrigger.Use(middleware.Auth)
		}
		apiTrigger.Any("/cron", v1.CronTaskController.TriggerTask)
	}
}
