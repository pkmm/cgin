package model

type SyncDetail struct {
	Id            uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	StudentId     uint64 `json:"student_id" gorm:"index"`
	StudentNumber string `json:"student_number" gorm:"unique;type:varchar(64);not null;"`
	Count         int    `json:"count" gorm:"default:0;"`
	CostTime      string `json:"cost_time" gorm:"size:10;default:null"`
	Info          string `json:"info" gorm:"size:255;default:null"`
	Model
}
