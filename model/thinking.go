package model

// this is model, which to record some interesting sentence.
// what is significance of living.

type Thinking struct {
	Id        uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId    uint64 `json:"user_id" gorm:"index:user_id;not null;"`
	Content   string `json:"content" gorm:"type:varchar(255);not null;"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
	Model
}
