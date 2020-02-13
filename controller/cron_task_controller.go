package controller

import (
	"cgin/conf"
	"cgin/controller/context_helper"
	"cgin/errno"
	"cgin/task"
	"github.com/gin-gonic/gin"
)

type cronTaskController struct {
}

var CronTaskController = cronTaskController{}

// @Summary 定时任务触发器
// @Produce json
// @Param job_name query string true "job_name" Enums(sign_baidu_tieba, not_exist)
// @Router /trigger/cron [get]
// @Failure 200 {object} service.Response
// @Success 200 {object} service.Response
func (this *cronTaskController) TriggerTask(c *gin.Context) {
	helper := context_helper.New(c)
	if conf.AppConfig.String("appEnv") != "dev" {
		helper.GetAuthUserId() // prod 需要认证
		// TODO 记录操作人
	}
	// 非生产环境
	jobName := ""
	jobName = helper.GetString("job_name")
	flag := false
	for _, taskName := range task.Tasks {
		if jobName == taskName {
			flag = true
			break
		}
	}
	if !flag {
		panic(errno.NormalException.ReplaceErrorMsgWith("未找到指定的任务"))
	}
	go func() {
		switch jobName {
		case task.FlagBaiduTiebaSign:
			task.SignBaiduForums()
		case task.FlagSyncStudentScore:
			task.UpdateStudentScore()
		default:
			return
		}
	}()
	helper.Response("任务已经在后台执行，请稍后查看")
}
