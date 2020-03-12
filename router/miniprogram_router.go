package router

import (
	"cgin/controller"
	"cgin/middleware"
	"github.com/gin-gonic/gin"
)

func mapMiniProgramRouter(router *gin.Engine) {
	apiMiniProgram := router.Group(MiniProgram)
	{
		// 不需要认证的
		apiMiniProgram.POST("/get_index_preference", controller.MiniProgramController.GetIndexPreference)
		apiMiniProgram.POST("/set_index_config", controller.MiniProgramController.SetIndexConfig)
		apiMiniProgram.GET("/get_notifications", controller.MiniProgramController.GetNotification)
		apiMiniProgram.GET("/get_sponsors", controller.MiniProgramController.GetSponsors)

		// 以下的API 需要认证
		apiMiniProgramNeedAuth := apiMiniProgram.Use(middleware.Auth)
		apiMiniProgramNeedAuth.POST("/config_menu", controller.MiniProgramController.DisposeMenu)
		apiMiniProgramNeedAuth.POST("/change_notification", controller.MiniProgramController.UpdateOrCreateNotification)
		apiMiniProgramNeedAuth.GET("/get_hermann_memorial", controller.HermannRememberController.GetTodayTaskInfo)
		apiMiniProgramNeedAuth.POST("/add_hermann_memorial", controller.HermannRememberController.SaveUserRememberTask)
	}
}
