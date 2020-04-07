package model

import (
	"cgin/conf"
)

// 小程序页面上半部显示的一些配置

type IndexConfig struct {
	Id         uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	Slogan     string `json:"slogan" gorm:"type:varchar(64);default:null;"`
	Motto      string `json:"motto" gorm:"type:varchar(64);default:null;"`
	ImageUrl   string `json:"image_url" gorm:"type:varchar(255);default:null;"`
	ImageStyle string `json:"image_style" gorm:"type:varchar(255);default:null;"`
	Disabled   bool   `json:"disabled" gorm:"default:false;"`
	Model
}

func (i *IndexConfig) Save() (err error, indexConfig *IndexConfig) {
	err = conf.DB.Create(i).Error
	return err, i
}

func (i *IndexConfig) GetLatest() (error, *IndexConfig) {
	var ic IndexConfig
	err := conf.DB.Last(&ic).Error
	return err, &ic
}
