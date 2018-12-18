package task

// 后台运行的任务
import (
	"database/sql"
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
	"pkmm_gin/model"
	"pkmm_gin/utility/zf"
	"sync"
	"time"
)

const (
	TEST               = "test"
	SYNC_STUDENT_SCORE = "sync student score"
	RUNNING            = 1 // 任务正在运行
	END                = 2 // 任务结束运行
)

// === 线程安全的map ===
type SafeMap struct {
	sync.RWMutex
	Map map[string]int
}

func newSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[string]int)
	return sm
}

func (this *SafeMap) readSafeMap(key string) (int, bool) {
	this.RLock()
	value, ok := this.Map[key]
	this.RUnlock()
	return value, ok
}

func (this *SafeMap) writeSafeMap(key string, value int) {
	this.Lock()
	this.Map[key] = value
	this.Unlock()
}

// === 线程安全的 map End ===

// 同步学生成绩 消费者线程产生的结果数据
type SyncStudentScoreResult struct {
	model.Student
	UseTime   string // 同步耗时
	SyncCount int    // 同步下来的课程数量
}

func init() {
	// 分配每一个goroutine 一个id, 避免任务的重叠运行
	runningTask := newSafeMap()

	c := cron.New()
	// 定义任务列表

	// 登录api, 移除过期的 nonce
	c.AddFunc("*/60 * * * * *", func() {
		ts := time.Now().Add(-1 * time.Minute).Unix()
		myRedis := model.GetRedis()
		defer myRedis.Close()
		myRedis.Do("ZREMRANGEBYSCORE", "api_collection_ttl", 0, ts)
		logs.Error("remove expired nonce.")
	})

	// 同步学生的成绩
	c.AddFunc("0 */10 * * * *", func() {
		if v, ok := runningTask.readSafeMap(SYNC_STUDENT_SCORE); ok && v == RUNNING {
			return
		}
		runningTask.writeSafeMap(SYNC_STUDENT_SCORE, RUNNING)

		var (
			id       int64
			num      string
			pwd      string
			rows     *sql.Rows
			err      error
			total    int64 = 0
			offset   int64 = 0
			page     int64 = 1
			pageSize int64 = 5
		)
		total, err = model.CanSyncCount()
		if err != nil {
			logs.Error(err.Error())
			return
		}

		goroutine := 10
		reqCh := make(chan model.Student, pageSize)
		resCh := make(chan SyncStudentScoreResult, goroutine<<1)
		closeCh := make(chan int)

		// 消费者
		for i := 0; i < goroutine; i++ {
			go func() {
				for {
					stu := <-reqCh
					c := zf.NewCrawl(stu.Num, stu.Pwd)
					beginTime := time.Now()
					retry := 5
					var scores []*zf.Score
					var err error
					for i := 0; i < retry; i++ {
						scores, err = c.LoginScorePage()
						if err == nil {
							break
						} else {
							// 如果是密码错误，当天就不再同步信息了，因为密码错误五次会锁定账号一天
							if err.Error() == zf.LOGIN_ERROR_MSG_WRONG_PASSWORD || err.Error() == zf.LOGIN_ERROR_MSG_NOT_VALID_USER {
								err := model.UpdateUserCanSync(stu.Id, 0)
								logs.Error(err)
								break
							}
						}
					}
					// 处理成绩
					if len(scores) != 0 {
						for _, s := range scores {
							score := &model.Score{
								Xn:        s.Xn,
								Xq:        s.Xq,
								Kcmc:      s.Kcmc,
								Cj:        s.Cj,
								Jd:        s.Jd,
								Cxcj:      s.Cxcj,
								Bkcj:      s.Bkcj,
								Xf:        s.Xf,
								Type:      s.Type,
								StudentId: stu.Id,
								UpdatedAt: time.Now(),
							}
							model.UpdateOrCreateScore(score)
						}
					}
					useTime := time.Since(beginTime)
					resCh <- SyncStudentScoreResult{Student: stu, UseTime: useTime.String(), SyncCount: len(scores)}
				}
			}()
		}

		// 生产者 分块查询，一直生产
		go func() {
			totalPage := (total + pageSize - 1) / pageSize
			for ; page <= totalPage; page++ {
				offset = (page - 1) * pageSize

				rows, err = model.GetDB().Raw(
					"SELECT id, num, pwd FROM students JOIN (SELECT id FROM students WHERE can_sync = 1 ORDER BY id ASC LIMIT ?, ?) AS stu_tmp USING(id)",
					offset,
					pageSize,
				).Rows()
				if err != nil {
					logs.Error(err.Error() + "???")
					continue
				}
				cnt := 0
				for rows.Next() {
					rows.Scan(&id, &num, &pwd)
					cnt++
					reqCh <- model.Student{Num: num, Pwd: pwd, Id: id}
				}
				// 没有成绩的话，结束 防止阻塞
				if cnt == 0 {
					close(closeCh)
					return
				}
			}
			defer rows.Close()
		}()

		// 读取结果
		// 全部学生的result
		// 此处的total不要算错，否则会阻塞，可以加一个超时
		for i := int64(0); i < total; i++ {
			select {
			case result := <-resCh:
				model.UpdateSyncDetail(model.SyncDetail{
					StuNo:     result.Num,
					LessonCnt: result.SyncCount,
					CostTime:  result.UseTime,
					UpdatedAt: time.Now(),
				})
			case <-closeCh:
				// TODO.
				logs.Info("no users need sync.")
				break
			case <-time.After(time.Duration(600 * time.Second)): // 600s超时的时间
				logs.Error("同步学生成绩 发生阻塞, 超时结束")
				break
			}
		}

		// 本次同步结束 清理操作
		runningTask.writeSafeMap(SYNC_STUDENT_SCORE, END)

	})

	// 在每天即将结束的时候，复位user的can_sync字段
	c.AddFunc("0 55 11 * * *", func() {
		model.SetCanSync(1)
	})

	c.AddFunc("0 */10 * * * *", func() {
		logs.Info(time.Now().Format("2006-01-02 15:04:05"))
	})
	//  开始任务
	c.Start()
}
