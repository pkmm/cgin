package service

import (
	"cgin/model"
	"sync"
)

// User Service
var User = &userService{
	mutex: &sync.Mutex{},
}

type userService struct {
	mutex *sync.Mutex
}

func (serv *userService) GetUser(id uint64) *model.User {
	user := &model.User{}
	if err := db.Where("`id` = ?", id).First(&user).Error; err != nil {
		// TODO: log
		return nil
	}
	return user
}

func (serv *userService) GetUserByOpenId(openId string) *model.User {
	user := &model.User{}
	if err := db.Model(&model.User{}).Where("`open_id` = ? ", openId).Preload("Student").First(&user).Error; err != nil {
		// TODO: log
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

func (serv *userService) CreateUserWithOpenId(openId string) *model.User {
	user := &model.User{
		OpenId: openId,
	}
	if err := db.Save(user).Error; err != nil {
		return nil
	}
	return user
}

func (serv *userService) GetStudentByUserId(userId uint64) *model.Student {
	student := &model.Student{}
	if err := db.Where("user_id = ?", userId).First(&student); err != nil {
		return nil
	}
	return student
}

func (serv *userService) UpdateStudentInfoByUserId(studentNumber, password string, userId uint64) *model.Student {
	student := &model.Student{}
	if err := db.Where("user_id = ?", userId).
		Assign(model.Student{Number: studentNumber, Password: password, UserId: userId}).
		FirstOrCreate(student).Error; err != nil {
		return nil
	}
	return student
}
