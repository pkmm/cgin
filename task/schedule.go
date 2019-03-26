package task

// 后台运行的任务
import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
	"cgin/zcmuES"
	"github.com/robfig/cron"
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

func init() {

	runningTask = newSafeMap()

	c := cron.New()
	// 定义任务列表

	// 同步学生的成绩
	c.AddFunc("0 */10 * * * *", taskWrapper(updateStudentScore, flagSyncStudentScore))

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 11 * * *", func() {
		service.User.ResetSyncStatus()
	})

	c.AddFunc("0 */10 * * * *", func() {
		conf.AppLogger.Info(time.Now().Format("2006-01-02 15:04:05"))
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

func updateStudentScore() {
	startAt := time.Now()
	// 执行的线程的数量
	workerCount := 10
	// 队列的大小
	queueSize := 32
	queue := make(chan *model.User, queueSize)

	// chunk todo
	users := service.User.GetCanSyncUsers(0, 100000)
	jobsSize := len(users)
	jobs := make(chan struct{}, jobsSize) // 一共有多少的人需要同步，阻塞至结束
	go func(users []*model.User) {
		for _, user := range users {
			queue <- user
		}
		close(queue)
	}(users)

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				var err error
				user, ok := <-queue
				if !ok {
					return // 直接结束这个for循环
				}
				beginAt := time.Now()
				conf.AppLogger.Info("begin sync student[num: %s] scores.", user.Num)
				zfWorker, err := zcmuES.NewCrawl(user.Num, user.Pwd)
				if err != nil {
					conf.AppLogger.Error("init crawl for user[num: %s] failed.", user.Num)
					jobs <- struct{}{}
					continue // 继续执行下一位的任务
				}
				// retry when err is verify code wrong.
				retry := 3
				var scores []*zcmuES.Score
				for tryTimes := 0; tryTimes <= retry; tryTimes++ {
					scores, err = zfWorker.GetScores()
					if err == nil {
						break
					} else if !zfWorker.CanContinue() {
						break
					}
				}

				if err != nil {
					conf.AppLogger.Error("sync student[num: %s] scores failed. reason: %s", user.Num, err.Error())
					service.User.SetUserAutoSyncStatus(user.ID, 0)
					jobs <- struct{}{}
					continue
				}

				if uint64(len(scores)) == service.ScoreService.GetUserScoreCount(user.ID) {
					// 成绩数量没有发生变化，那么就不进行同步了
					conf.AppLogger.Info("student[num: %s] scores not changed. current count[count: %d]", user.Num, len(scores))
					jobs <- struct{}{}
					continue
				} else {
					modelScores := make([]*model.Score, 0)
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
						modelScores = append(modelScores, score)
					}
					service.ScoreService.BatchCreate(modelScores)
				}

				endAt := time.Since(beginAt)
				service.ScoreService.UpdateSyncDetail(&model.SyncDetail{
					StuNo:        user.Num,
					LessonCnt:    len(scores),
					CostTime:     endAt.String(),
					FailedReason: "",
				})
				jobs <- struct{}{}
			}
		}()
	}

	tmp := jobsSize
	for tmp > 0 {
		<-jobs // 阻塞等待所有线程执行任务结束
		tmp--
	}
	stopAt := time.Since(startAt)
	conf.AppLogger.Info("sync %d students scores finish, use time %s", jobsSize, stopAt.String())
}
