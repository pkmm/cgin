package model

import "cgin/constant"

// 小程序的首页菜单的配置项
type Menu struct {
	Id          uint64                             `json:"id" gorm:"primary_key;auto_increment;"`
	Desp        string                             `json:"desp" gorm:"type:varchar(64);default:NULL;"`
	Title       string                             `json:"title" gorm:"type:varchar(16);NOT NULL;unique;"`
	Icon        string                             `json:"icon" gorm:"type:varchar(64);default:NULL;"`
	ActionType  constant.MiniProgramMenuActionType `json:"action_type" gorm:"default:0;"`
	ActionValue string                             `json:"action_value" gorm:"type:varchar(64);default:NULL;"`
	Disabled    bool                               `json:"disabled" gorm:"default:false;"`
	Model
}
