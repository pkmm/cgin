package controller

import (
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type dailyController struct {
	BaseController
}

var DailyController = &dailyController{}

func (d *dailyController) GetImage(c *gin.Context) {
	d.Response(c, gin.H{
		"image": service.DailyService.GetImage(),
	})
}

func (d *dailyController) GetSentence(c *gin.Context) {
	d.Response(c, gin.H{
		"sentence": service.DailyService.GetSentence(),
	})
}
