package system

import (
	"cgin/global"
	"cgin/model/system"
	"errors"
	"github.com/jinzhu/gorm"
)

type ApiService struct {
}

var ApiServiceApp = new(ApiService)

func (a *ApiService) CreateUser(user system.DeliUser) (err error) {
	if !errors.Is(global.DB.Where("username = ?", user.Username).First(&system.DeliUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户已经存在!")
	}
	return global.DB.Create(&user).Error
}

func (a *ApiService) GetUser(username string) (user *system.DeliUser, err error) {
	err = global.DB.Where("username = ?", username).First(&user).Error
	return user, err
}
