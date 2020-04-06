package model

import (
	"cgin/conf"
	"cgin/util"
	"github.com/jinzhu/gorm"
)

type ImageStorageInfo struct {
	Id  uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	Url string `json:"url"`
	Model
}

func (i *ImageStorageInfo) Create() (*ImageStorageInfo, error) {
	err := conf.DB.Where("date(created_at) = ?", util.Date()).Assign(*i).FirstOrCreate(i).Error
	return i, err
}

func (i *ImageStorageInfo) FindTodayImage() (_ *ImageStorageInfo, isFound bool) {
	date := util.Date()
	var ret ImageStorageInfo
	err := conf.DB.Where("Date(created_at) = ?", date).First(&ret).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false
	}
	if err != nil {
		return nil, true
	}
	return &ret, true
}
