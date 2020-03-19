package model

import "cgin/model/modelInterface"

// this is model, which to record some interesting sentence.
// what is significance of living.

type Thinking struct {
	Id        uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId    uint64 `json:"user_id" gorm:"index:user_id;not null;"`
	Content   string `json:"content" gorm:"type:varchar(255);not null;"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
	Model
}

func (t *Thinking) GetList(info modelInterface.PageSizeInfo) (err error, data interface{}, total int) {
	err, query, total := basicPagination(info, t)
	if err != nil {
		return err, nil, 0
	} else {
		var result []*Thinking
		err = query.Find(&result).Error
		return err, result, total
	}
}
