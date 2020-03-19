package task

import (
	"cgin/conf"
	"cgin/model"
	"cgin/service/workerpool"
	"cgin/util"
	"cgin/zcmu"
)

func UpdateStudentScore() {

	_updateStudentScore()
	return
	//
	//startAt := time.Now()
	//// 执行的线程的数量
	//workerCount := 10
	//
	//// chunk todo
	//students, err := service.StudentService.GetStudentNeedSyncScore(0, 100000)
	//if err != nil {
	//	conf.Logger.Error("Get student for sync student scores failed %#v", err)
	//	return
	//}
	//if len(students) == 0 {
	//	conf.Logger.Info("No students need sync.")
	//	return
	//}
	//// 队列的大小
	//queueSize := 32
	//if queueSize > len(students) {
	//	queueSize = len(students)
	//}
	//queue := make(chan *model.Student, queueSize)
	//outputQueue := make(chan *model.SyncDetail, queueSize)
	//
	//// 生产者 产生任务数据
	//go func(users []*model.Student) {
	//	for _, student := range students {
	//		queue <- student
	//	}
	//	close(queue)
	//}(students)
	//
	//for i := 0; i < workerCount; i++ {
	//	go func() {
	//		for {
	//			var err error
	//			student, ok := <-queue
	//			if !ok {
	//				return // 已经没有任务可以获取，结束工作线程
	//			}
	//			beginAt := time.Now()
	//			output := &model.SyncDetail{
	//				StudentId:     student.Id,
	//				StudentNumber: student.Number,
	//			}
	//			conf.Logger.Info("begin sync student[num: %s] scores.", student.Number)
	//			zfWorker, err := zcmu.NewCrawl(student.Number, student.Password)
	//			if err != nil {
	//				conf.Logger.Error("init crawl for user[num: %s] failed.", student.Number)
	//				output.Info = err.Error()
	//				outputQueue <- output
	//				continue // 继续执行下一位的任务
	//			}
	//			// retry if err is verify code wrong.
	//			retry := 3
	//			var scores []*zcmu.Score
	//			for tryTimes := 0; tryTimes <= retry; tryTimes++ {
	//				scores, err = zfWorker.GetScores()
	//				if err == nil {
	//					break
	//				} else if !zfWorker.CanContinue() {
	//					break
	//				}
	//			}
	//
	//			if err != nil {
	//				output.Info = err.Error()
	//				conf.Logger.Error("sync student[num: %s] scores failed. reason: %s", student.Number, err.Error())
	//				err = service.StudentService.UpdateStudentSyncStatus(student.Id, false)
	//				outputQueue <- output
	//				continue
	//			}
	//
	//			if uint64(len(scores)) == service.StudentService.GetScoreCount(student.Id) {
	//				// 成绩数量没有发生变化, 按照算法随机尝试更新
	//				x := rand.Intn(100)
	//				if x > 3 {
	//					conf.Logger.Info("student[num: %s] scores not changed. current count[count: %d]", student.Number, len(scores))
	//					outputQueue <- output
	//					continue
	//				}
	//			}
	//
	//			// 更新成绩
	//			modelScores := make([]*model.Score, 0)
	//			for _, s := range scores {
	//				score := &model.Score{}
	//				util.BeanDeepCopy(s, score)
	//				score.StudentId = student.Id
	//				modelScores = append(modelScores, score)
	//			}
	//			service.ScoreService.BatchCreate(modelScores)
	//
	//			endAt := time.Since(beginAt)
	//			output.count = len(scores)
	//			output.CostTime = endAt.String()
	//			outputQueue <- output
	//		}
	//	}()
	//}
	//
	//for i := 0; i < len(students); i++ {
	//	service.StudentService.UpdateSyncDetail(<-outputQueue)
	//}
	//
	//stopAt := time.Since(startAt)
	//conf.Logger.Info("sync %d students scores finish, use time %s", len(students), stopAt.String())
}

func _updateStudentScore() {
	// TODO: chunk 分块查询 优化查询 保存成绩， 更新爬虫的代码
	err, students, _ := new(model.Student).GetStudentsNeedSyncScore(1, 10)
	if err != nil {
		conf.Logger.Info("查询学生成绩失败 [%s]", err.Error())
		return
	}
	ts := make([]*workerpool.Task, len(*students))
	for i, stu := range *students {
		stuCopy := stu
		ts[i] = workerpool.NewTask(func() error {
			wk, err := zcmu.NewCrawl(stuCopy.Number, stuCopy.Password)
			if err != nil {
				conf.Logger.Error("Init crawl of student[%s] failed [%s]", stuCopy.Number, err.Error())
			}
			scores, err := wk.GetScores()
			if err != nil {
				conf.Logger.Error("Get scores of student[%s] failed [%s]", stuCopy.Number, err.Error())
			}
			// 更新成绩
			modelScores := make([]*model.Score, 0, len(scores))
			for _, s := range scores {
				score := &model.Score{}
				util.BeanDeepCopy(s, score)
				score.StudentId = stuCopy.Id
				modelScores = append(modelScores, score)
			}
			model.BatchCreateScores(modelScores)
			return nil
		})
	}
	pool.AddTasks(ts)
}
