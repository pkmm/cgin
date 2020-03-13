package controller

import (
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/service"
	"cgin/util"
	"github.com/gin-gonic/gin"
)

type hermannRememberController struct {}

var HermannRememberController = &hermannRememberController{}

type taskDetail struct {
	Unit      uint          `json:"unit"`
	TotalUnit uint          `json:"total_unit"`
	StartAt   util.JSONTime `json:"start_at,string"`
}

// @Summary 背单词：今天的任务
// @Security ApiKeyAuth
// @Produce json
// @Router /mini_program/get_hermann_memorial [get]
// @Success 200 {object} service.Response
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

// @Summary 添加背单词的任务
// @Produce json
// @Security ApiKeyAuth
// @Router /mini_program/add_hermann_memorial [post]
// @Success 200 {object} service.Response
// @Param addData body co.AddHermannMemorial true "data"
func (h *hermannRememberController) SaveUserRememberTask(c *gin.Context) {
	helper := context_helper.New(c)
	var params taskDetail
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	err := service.HermannService.SaveTask(params.Unit, params.TotalUnit, params.StartAt, helper.GetAuthUserId())
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	newTask := service.HermannService.GetTaskRecord(helper.GetAuthUserId())
	helper.Response(newTask)
}
