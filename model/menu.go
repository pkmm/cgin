package model

import (
	"cgin/conf"
	"cgin/constant/miniprogram/menuaction"
	"errors"
	"github.com/jinzhu/gorm"
)

// 小程序的首页菜单的配置项
type Menu struct {
	Id          uint64                    `json:"id" gorm:"primary_key;auto_increment;"`
	Desp        string                    `json:"desp" gorm:"type:varchar(64);default:NULL;"`
	Title       string                    `json:"title" gorm:"type:varchar(16);NOT NULL;unique;"`
	Icon        string                    `json:"icon" gorm:"type:varchar(64);default:NULL;"`
	ActionType  menuaction.MenuActionType `json:"action_type" gorm:"default:0;"`
	ActionValue string                    `json:"action_value" gorm:"type:varchar(64);default:NULL;"`
	Disabled    bool                      `json:"disabled" gorm:"type:bool;"`
	Model
}

var (
	menuAlreadyExist = errors.New("菜单项已经存在")
)

func (m *Menu) CreateMenu() error {
	err, _ := m.getMenuByTitle(m.Title)
	if err == nil {
		return menuAlreadyExist
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return conf.DB.Create(m).Error
}

func (m *Menu) getMenuByTitle(title string) (err error, menu *Menu) {
	var _m Menu
	err = conf.DB.First(&_m, Menu{Title: title}).Error
	return err, &_m
}

func GetActiveMenus() (error, *[]*Menu) {
	var result []*Menu
	err := conf.DB.Find(&result, "disabled = ?", false).Error
	return err, &result
}
