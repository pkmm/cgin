package task

// 后台运行的任务
import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	"pkmm_gin/model"
	"pkmm_gin/service"
	"pkmm_gin/util/zf"
	"sync"
	"time"
)

const (
	flagSyncStudentScore = "sync student score"
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
var log *logs.BeeLogger

func init() {
	log = logs.NewLogger()
	log.SetLogger(logs.AdapterFile, `{"filename":"logs/pkmm_gin.log","level":7,"daily":true,"maxdays":2}`)
	log.EnableFuncCallDepth(true)
	log.Async(1e3)

	runningTask = newSafeMap()

	c := cron.New()
	// 定义任务列表

	// 同步学生的成绩
	c.AddFunc("*/1 * * * * *", taskWrapper(syncStudentScore, flagSyncStudentScore))

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 11 * * *", func() {
		service.User.ResetSyncStatus()
	})

	c.AddFunc("0 */10 * * * *", func() {
		log.Info(time.Now().Format("2006-01-02 15:04:05"))
	})
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
		runningTask.writeSafeMap(flagSyncStudentScore, END)
	}
}

func syncStudentScore() {

	// 同步学生成绩 消费者线程产生的结果数据
	type SyncStudentScoreResult struct {
		User         *model.User
		UseTime      string // 同步耗时
		SyncCount    int    // 同步下来的课程数量
		FailedReason string // 失败的原因
	}

	var (
		total    uint64 = 0
		offset   uint64 = 0
		page     uint64 = 1
		pageSize uint64 = 100
	)
	total = service.User.GetCanSyncCount()

	goroutine := 10
	reqCh := make(chan *model.User, pageSize)
	resCh := make(chan *SyncStudentScoreResult, goroutine<<1)
	closeCh := make(chan int)

	// 消费者
	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				user := <-reqCh
				worker, err := zf.NewCrawl(user.Num, user.Pwd)
				if err != nil {
					resCh <- &SyncStudentScoreResult{User: user, UseTime: "0", SyncCount: 0, FailedReason: err.Error()}
					log.Error("schedule: => new craw failed. " + err.Error())
					return
				}
				beginTime := time.Now()
				var usedTime time.Duration
				var syncCount = 0
				retry := 5
				//var scores []*zf.score
				for i := 0; i < retry; i++ {
					scores, err2 := worker.GetScores()
					err = err2
					log.Error(err.Error())
					syncCount = len(scores)
					if err == nil && syncCount != 0 {
						currentCount := service.ScoreService.GetUserScoreCount(user.ID)
						if currentCount == uint64(syncCount) { // 当前成绩的数量没有发生变化
							// todo.
						} else {
							for _, s := range scores {
								score := &model.Score{
									Xn:     s.Xn,
									Xq:     s.Xq,
									Kcmc:   s.Kcmc,
									Cj:     s.Cj,
									Jd:     s.Jd,
									Cxcj:   s.Cxcj,
									Bkcj:   s.Bkcj,
									Xf:     s.Xf,
									Type:   s.Type,
									UserId: user.ID,
								}
								service.ScoreService.UpdateOrCreateScore(score)
							}
						}
						usedTime = time.Since(beginTime)
						resCh <- &SyncStudentScoreResult{User: user,
							UseTime: usedTime.String(), SyncCount: syncCount}
					}
					// 如果是密码错误，当天就不再同步信息了，因为密码错误五次会锁定账号一天
					if worker.CanContinue() == false {
						service.User.SetUserAutoSyncStatus(user.ID, 0)
						usedTime = time.Since(beginTime)
						resCh <- &SyncStudentScoreResult{User: user,
							UseTime: usedTime.String(), SyncCount: syncCount, FailedReason: err.Error()}
						return
					}
				}
				resCh <- &SyncStudentScoreResult{User: user,
					UseTime: time.Since(beginTime).String(), SyncCount: 0, FailedReason: err.Error()}
			}
		}()
	}

	// 生产者 分块查询，一直生产
	go func() {
		totalPage := (total + pageSize - 1) / pageSize
		for ; page <= totalPage; page++ {
			offset = (page - 1) * pageSize
			users := service.User.GetCanSyncUsers(offset, pageSize)
			// 没有成绩的话，结束 防止阻塞
			if len(users) == 0 {
				close(closeCh)
				return
			}
			for _, user := range users {
				reqCh <- user
			}
		}
	}()

	// 读取结果
	// 全部学生的result
	// 此处的total不要算错，否则会阻塞，可以加一个超时
	for i := uint64(0); i < total; i++ {
		select {
		case result := <-resCh:
			service.ScoreService.UpdateSyncDetail(&model.SyncDetail{
				StuNo:        result.User.Num,
				LessonCnt:    result.SyncCount,
				CostTime:     result.UseTime,
				FailedReason: result.FailedReason,
			})
		case <-closeCh:
			// TODO.
			log.Info("no users need sync.")
			break
		case <-time.After(time.Duration(600 * time.Second)): // 600s超时的时间
			log.Error("同步学生成绩 发生阻塞, 超时结束")
			break
		}
	}

}
