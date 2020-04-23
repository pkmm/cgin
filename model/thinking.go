package model

// this is model, which to record some interesting sentence.
// what is significance of living.

type Thinking struct {
	Id      uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId  uint64 `json:"user_id" gorm:"index:user_id;not null;"`
	Content string `json:"content" gorm:"type:varchar(255);not null;"`
	From    string `json:"from" grom:"type:varchar(64);default null;"`
	Model
}

func (t *Thinking) GetList(page, size int) (error, interface{}, int) {
	err, query, total := basicPagination(page, size, t)
	if err != nil {
		return err, nil, 0
	} else {
		var result []*Thinking
		err = query.Find(&result).Error
		return err, result, total
	}
}
