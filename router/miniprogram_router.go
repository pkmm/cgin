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
		apiMiniProgram.GET("/get_index_preference", v1.MiniProgramController.GetIndexPreference)
		apiMiniProgram.POST("/set_index_config", v1.MiniProgramController.SetIndexConfig)
		apiMiniProgram.GET("/get_notifications", v1.MiniProgramController.GetNotifications)
		apiMiniProgram.GET("/get_sponsors", v1.MiniProgramController.GetSponsors)
		apiMiniProgram.GET("/send_template_msg", v1.MiniProgramController.SendTemplateMsg)

		// 以下的API 需要认证
		apiMiniProgramNeedAuth := apiMiniProgram.Use(middleware.Auth)
		apiMiniProgramNeedAuth.POST("/config_menu", v1.MiniProgramController.DisposeMenu)
		apiMiniProgramNeedAuth.POST("/change_notification", v1.MiniProgramController.UpdateOrCreateNotification)
		apiMiniProgramNeedAuth.GET("/get_hermann_memorial", v1.HermannRememberController.GetTodayTaskInfo)
		apiMiniProgramNeedAuth.POST("/add_hermann_memorial", v1.HermannRememberController.SaveUserRememberTask)
	}
}
