package model

type SyncDetail struct {
	ID           uint64 `json:"id" gorm:"primary_key"`
	StuNo        string `json:"stu_no" gorm:"column:stu_no;unique"`
	LessonCnt    int    `json:"lesson_cnt"`
	CostTime     string `json:"cost_time" gorm:"size:40"`
	FailedReason string `json:"failed_reason" gorm:"size:255;default:null"`
	Model
}
