package v1

import (
	"cgin/controller/co"
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/model"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type hermannRememberController struct{}

var HermannRememberController = &hermannRememberController{}

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
	var params co.AddHermannMemorial
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	dbModel := model.HermannMemorial{
		TotalUnit:    params.TotalUnit,
		StartAt:      params.StartAt,
		RememberUnit: params.Unit,
		UserId:       helper.GetAuthUserId(),
	}

	err := dbModel.UpdateOrCreate()
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	_, newTask := dbModel.GetOwnerTaskRecord()
	helper.Response(newTask)
}
