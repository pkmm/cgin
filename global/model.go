package global

import (
	"time"

	"gorm.io/gorm"
)

type GModel struct {
	ID        uint `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除的时间
}
