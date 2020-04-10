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

// @Summary 系统健康检查
// @Description 查看系统的内存信息
// @Router /health_check/mem [get]
// @Tags 系统信息
// @Success 200 object service.Response
func (h *healthController) MemoryInfo(c *gin.Context) {
	helper := contextHelper.New(c)
	helper.Response(service.HealthService.MemoryUseInfo())
}
