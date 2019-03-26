package task

// 后台运行的任务
import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
	"cgin/zcmuES"
	"fmt"
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

func syncStudentScore() {

	startAt := time.Now()
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
	resCh := make(chan *SyncStudentScoreResult, goroutine<<2)
	closeCh := make(chan int)

	// 消费者
	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				user := <-reqCh
				conf.AppLogger.Debug("开始同步：" + user.Num)
				worker, err := zcmuES.NewCrawl(user.Num, user.Pwd)
				if err != nil {
					resCh <- &SyncStudentScoreResult{User: user, UseTime: "0", SyncCount: 0, FailedReason: err.Error()}
					conf.AppLogger.Error("schedule: => new craw failed. " + err.Error())
					return
				}
				beginTime := time.Now()
				var usedTime time.Duration
				var syncedCount = 0
				retry := 5
				for i := 0; i < retry; i++ {
					scores, err := worker.GetScores()
					if err != nil {
						conf.AppLogger.Error(fmt.Sprintf("同步学生: %d %s", user.ID, err.Error()))
						// 如果是密码错误，当天就不再同步信息了，因为密码错误五次会锁定账号一天
						if worker.CanContinue() == false {
							service.User.SetUserAutoSyncStatus(user.ID, 0)
							usedTime = time.Since(beginTime)
							resCh <- &SyncStudentScoreResult{User: user,
								UseTime: usedTime.String(), SyncCount: syncedCount, FailedReason: err.Error()}
							return
						}
						continue
					}
					syncedCount = len(scores)
					currentCount := service.ScoreService.GetUserScoreCount(user.ID)
					if 0 == currentCount {
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
						break
					}
					if uint64(syncedCount) != currentCount {
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
						break
					}
				}
				user.StudentName = worker.GetStudentName()
				service.User.UpdateUserName(user.StudentName, user.ID)
				resCh <- &SyncStudentScoreResult{User: user,
					UseTime: time.Since(beginTime).String(), SyncCount: syncedCount, FailedReason: ""}
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
			conf.AppLogger.Info("no users need sync.")
			break
		case <-time.After(time.Duration(500 * time.Second)): // 500s超时的时间
			conf.AppLogger.Error("同步学生成绩 发生阻塞, 超时结束")
			break
		}
	}

	conf.AppLogger.Info("%d 用户 花费时间, %s", total, time.Since(startAt).String())

}
