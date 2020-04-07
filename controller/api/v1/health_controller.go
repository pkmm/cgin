package v1

import (
	"cgin/controller/contextHelper"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

// 健康检查

type healthController struct {
}

var HealthCheckController healthController

func (h *healthController) MemoryInfo(c *gin.Context) {
	helper := contextHelper.New(c)
	helper.Response(service.HealthService.MemoryUseInfo())
}
