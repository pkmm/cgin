package model

import "time"

type Model struct {
	CreatedAt time.Time  `json:"created_at" gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//DeletedAt *time.Time `json:"-" gorm:"index:idx_deleted_at"`
}
