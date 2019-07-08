package model

import "cgin/util"

// 艾宾斯浩背诵单词的记忆曲线
// TODO: 支持多个记录
type HermannMemorial struct {
	Id     uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId uint64 `json:"user_id" gorm:"index"`
	// 设置一个单元的数量 默认是每天两个单元
	RememberUnit uint `json:"remember_unit" gorm:"default:2"`
	// 设置一共有多少单元
	TotalUnit uint `json:"total_unit" gorm:"default:24"`
	// 计算当前进行到第几天
	StartAt util.JSONTime `json:"start_at" gorm:"default:current_timestamp"`
	Model
}
