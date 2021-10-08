package system

import (
	v1 "cgin/api/v1"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
}

func (a *ApiRouter) InitApiRouter(r *gin.RouterGroup) {
	apiRouter := r.Group("api")
	var api = v1.ApiGroupApp.SystemApiGroup.SystemApi
	{
		apiRouter.GET("index", api.Index)
	}
}
