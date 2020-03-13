package v1

import (
	"cgin/controller/context_helper"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type dailyController struct {}

var DailyController = &dailyController{}

// @Summary 随机图片
// @Security ApiKeyAuth
// @Router /daily/image [get]
// @Produce image/jpeg
// @Success 200 {object} service.Response
func (d *dailyController) GetImage(c *gin.Context) {
	//c.File(service.DailyService.GetImage())
	helper := context_helper.New(c)
	helper.Response(gin.H{
		"image_url": service.DailyService.GetImage(),
	})
}

// 每日一言的数据
// @Summary 一句话
// @Security ApiKeyAuth
// @Produce json
// @Router /daily/sentence [get]
// @Success 200 {object} service.Response
func (d *dailyController) GetSentence(c *gin.Context) {
	helper := context_helper.New(c)
	helper.Response(gin.H{
		"sentence": service.DailyService.GetSentence(),
	})
}
