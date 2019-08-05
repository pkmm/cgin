package task

// 后台运行的任务
import (
	"cgin/conf"
	"cgin/service"
	"github.com/robfig/cron"
	"sync"
	"time"
)

const (
	FlagSyncStudentScore = "sync student score"
	FlagBaiduTiebaSign   = "baidu tieba sign"
	//RUNNING
	RUNNING = 1 // 任务正在运行
	//END
	END = 2 // 任务结束运行
)

// === 线程安全的map ===
type safeMap struct {
	sync.RWMutex
	Map map[string]int
}

func newSafeMap() *safeMap {
	sm := new(safeMap)
	sm.Map = make(map[string]int)
	return sm
}

func (this *safeMap) readSafeMap(key string) (int, bool) {
	this.RLock()
	value, ok := this.Map[key]
	this.RUnlock()
	return value, ok
}

func (this *safeMap) writeSafeMap(key string, value int) {
	this.Lock()
	this.Map[key] = value
	this.Unlock()
}

// === 线程安全的 map End ===

// 分配每一个goroutine 一个id, 避免任务的重叠运行
var runningTask *safeMap

func init() {

	runningTask = newSafeMap()

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
		// before task
		if v, ok := runningTask.readSafeMap(flag); ok && v == RUNNING {
			return
		}
		runningTask.writeSafeMap(flag, RUNNING)
		cmd()
		// after task clean up
		runningTask.writeSafeMap(flag, END)
	}
}

