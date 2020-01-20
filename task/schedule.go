package task

// 后台运行的任务
import (
	"cgin/conf"
	"cgin/service"
	"cgin/util"
	"github.com/robfig/cron"
	"time"
)

const (
	FlagSyncStudentScore = "sync_student_score"
	FlagBaiduTiebaSign   = "sign_baidu_tieba"
)

var Tasks = []string{FlagBaiduTiebaSign, FlagSyncStudentScore}



// 分配每一个goroutine 一个id, 避免任务的重叠运行
var runningTask *util.SafeMap

func init() {

	runningTask = util.NewSafeMap()

	c := cron.New()
	// 定义任务列表

	// 同步学生的成绩
	c.AddFunc("0 */10 * * * *", taskWrapper(UpdateStudentScore, FlagSyncStudentScore))

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 11 * * *", func() {
		service.StudentService.RestSyncStatus()
	})

	// 测试用
	c.AddFunc("0 */10 * * * *", func() {
		conf.AppLogger.Info(time.Now().Format("2006-01-02 15:04:05"))
	})

	// 百度贴吧签到
	c.AddFunc("0 0 0 * * *", taskWrapper(SignBaiduForums, FlagBaiduTiebaSign))

	//  开始任务
	c.Start()
}

func taskWrapper(cmd func(), flag string) func() {
	return func() {
		if _, ok := runningTask.ReadSafeMap(flag); ok {
			return
		}
		runningTask.WriteSafeMap(flag, 1)
		cmd()
		runningTask.DeleteKey(flag)
	}
}

