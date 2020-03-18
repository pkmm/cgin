package model

import (
	"cgin/util"
)

type Notification struct {
	Id       uint64        `json:"id" gorm:"primary_key;auto_increment;"`
	Content  string        `json:"content" gorm:"type:varchar(255);default:null"`
	Disabled bool          `json:"disabled" gorm:"default:false;"`
	StartAt  util.JSONTime `json:"start_at" gorm:"default:current_timestamp;"`
	EndAt    util.JSONTime `json:"end_at" gorm:"default:current_timestamp;"`
	Model
}