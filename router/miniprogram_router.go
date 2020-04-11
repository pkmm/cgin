package router

import (
	"cgin/controller/api/v1"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func initMiniProgramRouter(miniProgramRouter *gin.Engine) {
	apiMiniProgram := miniProgramRouter.Group(MiniProgram)
	{
		// 不需要认证的
		apiMiniProgram.GET("/index_preferences", v1.MiniProgramController.IndexPreference)
		apiMiniProgram.GET("/notifications", v1.MiniProgramController.GetNotifications)
		apiMiniProgram.GET("/sponsors", v1.MiniProgramController.GetSponsors)
		apiMiniProgram.GET("/send_template_msg", v1.MiniProgramController.SendTemplateMsg)

		// 以下的API 需要认证
		apiMiniProgramNeedAuth := apiMiniProgram.Use(middleware.Auth())
		apiMiniProgramNeedAuth.POST("/index_config", v1.MiniProgramController.CreateIndexConfig)
		apiMiniProgramNeedAuth.POST("/menus", v1.MiniProgramController.CreateMenus)
		apiMiniProgramNeedAuth.PUT("/notifications", v1.MiniProgramController.UpdateOrCreateNotification)
		apiMiniProgramNeedAuth.GET("/hermann_memorials", v1.HermannRememberController.IndexHermannMemorial)
		apiMiniProgramNeedAuth.POST("/hermann_memorials", v1.HermannRememberController.CreateHermannMemorial)
	}
}
