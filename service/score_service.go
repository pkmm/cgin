package service

import (
	"cgin/model"
	"github.com/jinzhu/gorm"
	"sync"
)

type scoreService struct {
	mutex *sync.Mutex
}

var ScoreService = &scoreService{
	mutex: &sync.Mutex{},
}

func (s *scoreService) GetOwnScores(userId uint64) (error, *[]model.Score) {
	err, student := model.GetStudentByUserId(userId)
	if err != nil {
		// 记录没有找到
		if err == gorm.ErrRecordNotFound {
			return nil, &[]model.Score{}
		}
		return err, nil
	}
	return nil, &student.Scores
}

//func (s *scoreService) SaveStudentScoresFromCrawl(scores *[]zcmu.KcItem, studentId uint64) []*model.Score {
//	dbScores := make([]*model.Score, 0, len(*scores))
//	for _, s := range *scores {
//		modelScore := &model.Score{}
//		util.BeanDeepCopy(s, modelScore)
//		modelScore.StudentId = studentId
//		dbScores = append(dbScores, modelScore)
//	}
//	workerpool.TaskPool.AddTasks([]*workerpool.Task{workerpool.NewTask(func() {
//		model.BatchCreateScores(dbScores)
//	})})
//	//go s.BatchCreate(dbScores)
//	return dbScores
//}
