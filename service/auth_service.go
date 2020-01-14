package service

import (
	"cgin/model"
	"sync"
)

type authService struct {
	mutex *sync.Mutex
}

var AuthService = &authService{
	mutex: &sync.Mutex{},
}

// 小程序的注册方式
func (a *authService) LoginFromMiniProgram(openid string) *model.User {
	user := &model.User{}
	// 第一个记录或者是创建记录
	if err := db.Model(&model.User{}).Where("open_id = ?", openid).Preload("Student").Attrs(model.User{OpenId:openid}).FirstOrCreate(&user).Error; err != nil {
		return nil
	}
	return user
}
