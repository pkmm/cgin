package task

// 后台运行的任务
import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
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

var pool *workerpool.SimplePool

func init() {

	// 协程 作为全部的任务的执行器
	// 需要停止 stop 函数 才能结束全部的协程
	//pool = service.NewSimplePool(30)
	pool = workerpool.TaskPool
	pool.RunPool()

	runningTask = util.NewSafeMap()

	// 秒 分 时
	c := cron.New()
	// 定义任务列表

	// 同步学生的成绩
	c.AddFunc("0 */10 * * * *", taskWrapper(UpdateStudentScore, FlagSyncStudentScore))

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 23 * * *", func() {
		model.ResetStudentSyncScoreStatus()
	})

	// 测试用
	c.AddFunc("0 */10 * * * *", func() {
		task := workerpool.NewTask(func() {
			conf.Logger.Info("测试task任务: %s", time.Now().Format("2006-01-02 15:04:05"))
		})
		pool.AddTasks([]*workerpool.Task{task})
	})

	// 百度贴吧签到
	c.AddFunc("0 0 0 * * *", taskWrapper(SignBaiduForums, FlagBaiduTiebaSign))

	// 自己加载每日一图保存在weibo图床
	c.AddFunc("0 0 0 * * *", func() {
		service.DailyService.GetImage()
		//file, _ := ioutil.ReadFile(imagePath)
		//service.NewWeiBoStorage(conf.WeiBoCookie()).UploadImage(file)
	})

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
