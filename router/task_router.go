package router

import (
	"cgin/conf"
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func initTaskRouter(taskRouter *gin.Engine) {
	apiTrigger := taskRouter.Group(Trigger)
	{
		if conf.IsProd() {
			apiTrigger.Use(middleware.Auth)
		}
		apiTrigger.Any("/cron", v1.CronTaskController.TriggerTask)
	}
}
