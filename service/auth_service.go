package service

import (
	"cgin/conf"
	"cgin/model"
	"cgin/util"
)

type authService struct{}

var AuthService = &authService{}

// 小程序的注册方式
func (a *authService) LoginFromMiniProgram(openid string) (error, *model.User) {
	user := model.GetUserByOpenId(openid)
	if user != nil {
		return nil, user
	}
	// user not found
	user = &model.User{OpenId: openid, Password: util.RandomString(10)}
	user = user.Create()
	// query user form db
	user = model.GetUserByOpenId(openid)
	return nil, user
}

func (a *authService) LoginFromWebBrowser(username, password string) (error, *model.User) {
	user := model.User{}
	openId := "wb_" + util.GUID()
	err := conf.DB.Where("username = ?", username).
		Preload("Student").
		Preload("Role").
		Attrs(model.User{OpenId: openId, RoleId: model.RoleUser, Username: username, Password: password}).
		FirstOrCreate(&user).Error
	return err, &user
}
