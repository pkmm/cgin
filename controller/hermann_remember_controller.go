package controller

import (
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type hermannRememberController struct {
}

var HermannRememberController = &hermannRememberController{}

type taskDetail struct {
	Unit      uint          `json:"unit"`
	TotalUnit uint          `json:"total_unit"`
	StartAt   util.JSONTime `json:"start_at,string"`
}

func (h *hermannRememberController) GetTodayTaskInfo(c *gin.Context) {
	helper := context_helper.New(c)
	tasks, err := service.HermannService.GetTodayTask(helper.GetAuthUserId())
	if err != nil {
		panic(err)
	}
	helper.Response(gin.H{
		"tasks": tasks,
	})
}

func (h *hermannRememberController) SaveUserRememberTask(c *gin.Context) {
	helper := context_helper.New(c)
	var params taskDetail
	if err := c.BindJSON(&params); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	err := service.HermannService.SaveTask(params.Unit, params.TotalUnit, params.StartAt, helper.GetAuthUserId())
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	helper.Response(nil)
}
