package service

import (
	"cgin/model"
	"cgin/service/workerpool"
	"cgin/util"
	"cgin/zcmu"
	"sync"
)

type scoreService struct {
	mutex *sync.Mutex
}

var ScoreService = &scoreService{
	mutex: &sync.Mutex{},
}

func (s *scoreService) GetOwnScores(userId uint64) (err error, _scores *[]*model.Score) {
	err, student := model.GetStudentByUserId(userId)
	if student == nil {
		return nil, nil
	}
	err, scores := model.GetScoresByStudentId(student.Id)
	return err, &scores
}

func (s *scoreService) SaveStudentScoresFromCrawl(scores []*zcmu.Score, studentId uint64) []*model.Score {
	dbScores := make([]*model.Score, 0, len(scores))
	for _, s := range scores {
		modelScore := &model.Score{}
		util.BeanDeepCopy(s, modelScore)
		modelScore.StudentId = studentId
		dbScores = append(dbScores, modelScore)
	}
	workerpool.TaskPool.AddTasks([]*workerpool.Task{workerpool.NewTask(func() error {
		model.BatchCreateScores(dbScores)
		return nil
	})})
	//go s.BatchCreate(dbScores)
	return dbScores
}
