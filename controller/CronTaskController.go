package controller

import (
	"cgin/conf"
	"cgin/errno"
	"cgin/task"
	"github.com/gin-gonic/gin"
)

type cronTaskController struct {
	BaseController
}

var CronTaskController = cronTaskController{}

func (this *cronTaskController) TriggerTask(c *gin.Context) {
	if conf.AppConfig.String("appEnv") != "dev" {
		this.getAuthUserId(c) // prod 需要认证
		// TODO 记录操作人
	}
	// 非生产环境
	this.processParams(c)
	if jobName, ok := this.Params["job_name"].(string); !ok {
		panic(errno.NormalException.AppendErrorMsg("参数解析错误"))
	} else {
		go func() {
			switch jobName {
			case task.FlagBaiduTiebaSign:
				task.SignBaiduForums()
			case task.FlagSyncStudentScore:
				task.UpdateStudentScore()
			default:
				panic(errno.NormalException.ReplaceErrorMsgWith("未找到指定的任务"))
			}
		}()
	}
	this.response(c, "任务已经在后台执行，请稍后查看")
}