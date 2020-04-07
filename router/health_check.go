package router

import (
	v1 "cgin/controller/api/v1"
	"github.com/gin-gonic/gin"
)

const path = RootApiPrefix + "/health_check"

func initHealthCheck(healthCheckRouter *gin.Engine) {
	api := healthCheckRouter.Group(path)
	{
		api.GET("mem", v1.HealthCheckController.MemoryInfo)
	}
}
