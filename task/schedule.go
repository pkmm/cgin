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
	FlagBackupMysql      = "backup_mysql"
)

var Tasks = []string{FlagBaiduTiebaSign, FlagSyncStudentScore}

// 分配每一个goroutine 一个id, 避免任务的重叠运行
var runningTask *util.SafeMap

var pool *service.SimplePool

func init() {

	// 协程 作为全部的任务的执行器
	// 需要停止 stop 函数 才能结束全部的协程
	//pool = service.NewSimplePool(30)
	pool = service.TaskPool
	pool.RunPool()

	runningTask = util.NewSafeMap()

	// 秒 分 时
	c := cron.New()
	// 定义任务列表

	// 同步学生的成绩
	c.AddFunc("0 */10 * * * *", taskWrapper(UpdateStudentScore, FlagSyncStudentScore))

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 23 * * *", func() {
		service.StudentService.RestSyncStatus()
	})

	// 测试用
	c.AddFunc("0 */10 * * * *", func() {
		task := service.NewTask(func() error {
			conf.Logger.Info(time.Now().Format("2006-01-02 15:04:05"))
			return nil
		})
		pool.AddTasks([]*service.Task{task})
	})

	// 百度贴吧签到
	c.AddFunc("0 0 0 * * *", taskWrapper(SignBaiduForums, FlagBaiduTiebaSign))

	// 数据库备份
	c.AddFunc("0 0 3 * * *", taskWrapper(backupMysql, FlagBackupMysql))

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

func CleanPool() {
	if pool != nil {
		pool.Stop()
	}
}