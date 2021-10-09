package system

import (
	"cgin/global"
	"cgin/model/system"
)

type BaiduService struct {
}

var BaiduServiceApp = new(BaiduService)

func (b *BaiduService) GetUsers() (err error, users []system.BaiduUser) {
	err = global.DB.Find(&users).Error
	return err, users
}
