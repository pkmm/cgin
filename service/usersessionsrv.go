package service

import (
	"errors"
	"pkmm_gin/model"
	"pkmm_gin/util"
	"sync"
)

type userSessionService struct {
	mutex *sync.Mutex
}

var UserSessionService = &userSessionService{
	mutex: &sync.Mutex{},
}

func (srv *userSessionService) GetUserSession(userId uint64) *model.UserSession {
	ret := &model.UserSession{}
	if err := db.Where("`user_id` = ? AND `active` = 1", userId).First(&ret).Error; err != nil {
		if db.RecordNotFound() {
			ret = srv.CreateUserSession(userId)
			if ret != nil {
				return ret
			}
		}
		return nil
	}

	return ret
}

func (sev *userSessionService) getSessionIgnoreActive(userId uint64) *model.UserSession {
	ret := &model.UserSession{}
	if err := db.Where("`user_id` = ?", userId).Error; err != nil {
		return nil
	}

	return ret
}

func (srv *userSessionService) CreateUserSession(userId uint64) *model.UserSession {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	ret := &model.UserSession{}
	token := util.GenerateToken(userId)
	ret.AccessToken = token
	ret.UserId = userId
	tx := db.Begin()
	if err := tx.Create(ret).Error; err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()

	return ret
}

func (srv *userSessionService) UpdateUserSession(userSession *model.UserSession) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	if 0 == userSession.UserId {
		return errors.New("user id is empty")
	}

	userSession.Active = true
	userSession.AccessToken = util.GenerateToken(userSession.UserId)
	tx := db.Begin()
	if err := tx.Save(userSession).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return err
}
