package controller

import (
	"cgin/controller/context_helper"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type dailyController struct {}

var DailyController = &dailyController{}

// 直接返回图片，而不是常见的API JSON数据接口
// @Summary 随机图片
// @Router /daily/image [get]
// @Produce image/jpeg
// @Success 200 {object} service.Response
func (d *dailyController) GetImage(c *gin.Context) {
	c.File(service.DailyService.GetImage())
}

// 每日一言的数据
// @Summary 一句话
// @Produce json
// @Router /daily/sentence [get]
// @Success 200 {object} service.Response
func (d *dailyController) GetSentence(c *gin.Context) {
	helper := context_helper.New(c)
	helper.Response(gin.H{
		"sentence": service.DailyService.GetSentence(),
	})
}
