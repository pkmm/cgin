package v1

import (
	"cgin/controller/contextHelper"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type dailyController struct{}

var DailyController = &dailyController{}

// @Summary 随机图片
// @Tags Daily
// @Security ApiKeyAuth
// @Router /daily/images [get]
// @Produce image/jpeg
// @Success 200 object service.Response
func (d *dailyController) GetImage(c *gin.Context) {
	helper := contextHelper.New(c)
	helper.Response(gin.H{
		"image_url": service.DailyService.GetImage(),
	})
}

// @Summary 一句话
// @Security ApiKeyAuth
// @Tags Daily
// @Router /daily/sentences [get]
// @Success 200 object service.Response
func (d *dailyController) GetSentence(c *gin.Context) {
	helper := contextHelper.New(c)
	helper.Response(gin.H{
		"sentence": service.DailyService.GetSentence(),
	})
}
