package task

import (
	"cgin/conf"
	"cgin/model"
	"cgin/service"
	"cgin/util"
	"cgin/zcmu"
	"math/rand"
	"time"
)

func UpdateStudentScore() {
	startAt := time.Now()
	// 执行的线程的数量
	workerCount := 10

	// chunk todo
	students, err := service.StudentService.GetStudentNeedSyncScore(0, 100000)
	if err != nil {
		conf.AppLogger.Error("Get student for sync student scores failed ", err.Error())
		return
	}
	if len(students) == 0 {
		conf.AppLogger.Info("No students need sync.")
		return
	}
	// 队列的大小
	queueSize := 32
	if queueSize > len(students) {
		queueSize = len(students)
	}
	queue := make(chan *model.Student, queueSize)
	outputQueue := make(chan *model.SyncDetail, queueSize)

	// 生产者 产生任务数据
	go func(users []*model.Student) {
		for _, student := range students {
			queue <- student
		}
		close(queue)
	}(students)

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				var err error
				student, ok := <-queue
				if !ok {
					return // 已经没有任务可以获取，结束工作线程
				}
				beginAt := time.Now()
				output := &model.SyncDetail{
					StudentId:     student.Id,
					StudentNumber: student.Number,
				}
				conf.AppLogger.Info("begin sync student[num: %s] scores.", student.Number)
				zfWorker, err := zcmu.NewCrawl(student.Number, student.Password)
				if err != nil {
					conf.AppLogger.Error("init crawl for user[num: %s] failed.", student.Number)
					output.Info = err.Error()
					outputQueue <- output
					continue // 继续执行下一位的任务
				}
				// retry if err is verify code wrong.
				retry := 3
				var scores []*zcmu.Score
				for tryTimes := 0; tryTimes <= retry; tryTimes++ {
					scores, err = zfWorker.GetScores()
					if err == nil {
						break
					} else if !zfWorker.CanContinue() {
						break
					}
				}

				if err != nil {
					output.Info = err.Error()
					conf.AppLogger.Error("sync student[num: %s] scores failed. reason: %s", student.Number, err.Error())
					err = service.StudentService.UpdateStudentSyncStatus(student.Id, false)
					outputQueue <- output
					continue
				}

				if uint64(len(scores)) == service.StudentService.GetScoreCount(student.Id) {
					// 成绩数量没有发生变化, 按照算法随机尝试更新
					x := rand.Intn(100)
					if x > 3 {
						conf.AppLogger.Info("student[num: %s] scores not changed. current count[count: %d]", student.Number, len(scores))
						outputQueue <- output
						continue
					}
				}

				// 更新成绩
				modelScores := make([]*model.Score, 0)
				for _, s := range scores {
					score := &model.Score{}
					util.BeanDeepCopy(s, score)
					score.StudentId = student.Id
					modelScores = append(modelScores, score)
				}
				service.ScoreService.BatchCreate(modelScores)

				endAt := time.Since(beginAt)
				output.Count = len(scores)
				output.CostTime = endAt.String()
				outputQueue <- output
			}
		}()
	}

	for i := 0; i < len(students); i++ {
		service.StudentService.UpdateSyncDetail(<-outputQueue)
	}

	stopAt := time.Since(startAt)
	conf.AppLogger.Info("sync %d students scores finish, use time %s", len(students), stopAt.String())
}
