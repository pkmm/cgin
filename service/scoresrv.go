package service

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"pkmm_gin/model"
	"strings"
	"sync"
)

type scoreService struct {
	mutex *sync.Mutex
}

var ScoreService = &scoreService{
	mutex: &sync.Mutex{},
}

func (serv *scoreService) BatchCreate(scores []*model.Score) {
	sql := bytes.Buffer{}
	sql.WriteString("INSERT IGNORE INTO scores(user_id, xn, xq, kcmc, type, xf, jd, cj, bkcj, cxcj) ")
	binds := make([]interface{}, 0)
	for i, score := range scores {
		//fmt.Printf("%#v", score)
		if i == 0 {
			sql.WriteString("VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		} else {
			sql.WriteString(" ,(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		}
		binds = append(binds, score.UserId, score.Xn, score.Xq, score.Kcmc,
			score.Type, score.Xf, score.Jd, score.Cj, score.Bkcj, score.Cxcj)
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
		logs.Error("update or create user score failed." + err.Error())
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
