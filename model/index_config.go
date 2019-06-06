package model

// 小程序首页的一些配置

type IndexConfig struct {
	Id           uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	Slogan       string `json:"slogan" gorm:"type:varchar(64);default:null;"`
	ImageUrl     string `json:"image_url" gorm:"type:varchar(255);default:null;"`
	ImageStyle   string `json:"image_style" gorm:"type:varchar(255);default:null;"`
	Disabled     bool   `json:"disabled" gorm:"default:false;"`
	Model
}
