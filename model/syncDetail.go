package model

import "time"

type SyncDetail struct {
	Id        int64     `json:"id"`
	StuNo     string    `json:"stu_no" gorm:"column:stu_no"`
	LessonCnt int       `json:"lesson_cnt"`
	CostTime  string    `json:"cost_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SyncDetail) TableName() string {
	return "sync_detail"
}

// updateOrCreate
func UpdateSyncDetail(syncDetail SyncDetail) SyncDetail {
	db.Where(SyncDetail{StuNo: syncDetail.StuNo}).Assign(syncDetail).FirstOrCreate(&syncDetail)
	return syncDetail
}
