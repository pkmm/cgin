package controller

import (
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type hermannRememberController struct {
	BaseController
}

var HermannRememberController = &hermannRememberController{}

type taskDetail struct {
	Unit      uint          `json:"unit"`
	TotalUnit uint          `json:"total_unit"`
	StartAt   util.JSONTime `json:"start_at,string"`
}

func (h *hermannRememberController) GetTodayTaskInfo(c *gin.Context) {
	h.getAuthUserId(c)
	tasks, err := service.HermannService.GetTodayTask(h.UserId)
	if err != nil {
		panic(err)
	}
	h.response(c, gin.H{
		"tasks": tasks,
	})
}

func (h *hermannRememberController) SaveUserRememberTask(c *gin.Context) {
	h.getAuthUserId(c)
	var params taskDetail
	if err := c.BindJSON(&params); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	err := service.HermannService.SaveTask(params.Unit, params.TotalUnit, params.StartAt, h.UserId)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	h.response(c, nil)
}
