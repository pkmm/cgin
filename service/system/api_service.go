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

func (a *ApiService) CreateUser(user system.User) (err error) {
	if !errors.Is(global.G_DB.Where("username = ?", user.Username).First(&system.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户已经存在!")
	}
	return global.G_DB.Create(&user).Error
}

func (a *ApiService) GetUser(username string) (user *system.User, err error) {
	err = global.G_DB.Where("username = ?", username).First(&user).Error
	return user, err
}
