package service

import (
	"bytes"
	"cgin/conf"
	"cgin/model"
	"strings"
	"sync"
	"time"
)

type scoreService struct {
	mutex *sync.Mutex
}

var ScoreService = &scoreService{
	mutex: &sync.Mutex{},
}

func (serv *scoreService) BatchCreate(scores []*model.Score) {
	if len(scores) == 0 {
		return
	}
	sql := bytes.Buffer{}
	sql.WriteString("INSERT IGNORE INTO scores(user_id, xn, xq, kcmc, type, xf, jd, cj, bkcj, cxcj, created_at, updated_at) ")
	binds := make([]interface{}, 0)
	for i, score := range scores {
		//fmt.Printf("%#v", score)
		if i == 0 {
			sql.WriteString("VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		} else {
			sql.WriteString(" ,(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		}
		binds = append(binds, score.UserId, score.Xn, score.Xq, score.Kcmc,
			score.Type, score.Xf, score.Jd, score.Cj, score.Bkcj, score.Cxcj,
			time.Now().Unix(), time.Now().Unix())
	}
	sqls := strings.Trim(sql.String(), ",")
	//fmt.Println(sqls, binds)
	//return
	db.Exec(sqls, binds...)
}

func (serv *scoreService) UpdateOrCreateScore(score *model.Score) *model.Score {
	if err := db.Where(&model.Score{UserId: score.UserId, Xn: score.Xn, Xq: score.Xq, Kcmc: score.Kcmc}).
		Assign(score).
		FirstOrCreate(&score).Error; err != nil {
		conf.AppLogger.Error("update or create user score failed." + err.Error())
		return nil
	}

	return score
}

func (serv *scoreService) UpdateSyncDetail(syncDetail *model.SyncDetail) *model.SyncDetail {
	if err := db.Where(&model.SyncDetail{StuNo: syncDetail.StuNo}).
		Assign(model.SyncDetail{CostTime: syncDetail.CostTime,
			FailedReason: syncDetail.FailedReason, LessonCnt: syncDetail.LessonCnt}).
		FirstOrCreate(&syncDetail).Error; err != nil {
		return nil
	}

	return syncDetail
}

func (serv *scoreService) GetUserScoreCount(userId uint64) (count uint64) {
	if err := db.Model(&model.Score{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0
	}

	return
}

func (serv *scoreService) GetOwnScores(userId uint64) (scores []*model.Score) {
	if err := db.Where(&model.Score{UserId: userId}).Find(&scores).Error; err != nil {
		return
	}

	return scores
}
