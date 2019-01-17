package service

import (
	"github.com/astaxie/beego/logs"
	"pkmm_gin/model"
	"sync"
)

// User Service
var User = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
}

func (serv *userService) CheckAndGetUserByUserIdAndAccessToken(userId uint64, accessToken string) *model.User {
	userSession := &model.UserSession{}
	if err := db.Where("`user_id` = ? AND `access_token` = ?", userId, accessToken).
		First(userSession).Error; err != nil {
		// todo log
		return nil
	}

	return serv.GetUser(userSession.UserId)
}

func (serv *userService) GetUser(id uint64) *model.User {
	user := &model.User{}
	if err := db.Where("`id` = ?", id).First(user).Error; err != nil {
		// log todo
		return nil
	}

	return user
}

func (serv *userService) GetUserByOpenId(openId string) *model.User {
	user := &model.User{}
	if err := db.Where("`open_id` = ? ", openId).First(&user).Error; err != nil {
		// todo
		return nil
	}

	return user
}

func (serv *userService) UpdateUser(user *model.User) error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	tx := db.Begin()
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (serv *userService) GetCanSyncCount() uint64 {
	count := 0
	if err := db.Model(&model.User{}).Where("`can_sync` = 1").Count(&count).Error; err != nil {
		return 0
	}

	return uint64(count)
}

func (serv *userService) SetUserAutoSyncStatus(userId uint64, canSync int) {
	db.Model(&model.User{}).Where("`id` = ?", userId).Update("can_sync", canSync)
}

func (serv *userService) ResetSyncStatus() {
	db.Model(&model.User{}).Updates(map[string]interface{}{"can_sync": 1})
}

func (serv *userService) GetCanSyncUsers(offset, limit uint64) (users []*model.User) {
	if err := db.Model(&model.User{}).Where("`can_sync` = 1").
		Order("`id` DESC").
		Limit(limit).
		Offset(offset).Find(&users).Error; err != nil {
		logs.Error("get sync users failed" + err.Error())
	}

	return users
}
