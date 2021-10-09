package task

// 后台运行的任务
import (
	"cgin/service/workerpool"
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
var pool *workerpool.Pool

func init() {

	pool, err := workerpool.NewPool(20, time.Second*10)
	if err != nil {
		panic("initialize worker pool failed.")
	}
	defer pool.Close()

	runningTask = util.NewSafeMap()

	// 秒 分 时
	c := cron.New()
	// 定义任务列表

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
