package model

import (
	"cgin/util"
)

type Model struct {
	CreatedAt util.JSONTime  `json:"created_at" gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt util.JSONTime  `json:"updated_at" gorm:"DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//DeletedAt *time.Time `json:"-" gorm:"index:idx_deleted_at"`
}
