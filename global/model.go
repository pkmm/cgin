package global

import (
	"gorm.io/gorm"
	"time"
)

type GModel struct {
	ID        uint `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除的时间
}
