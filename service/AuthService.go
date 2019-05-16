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
	if user := User.GetUserByOpenId(openid); user != nil {
		return user
	}
	return User.CreateUserWithOpenId(openid)
}
