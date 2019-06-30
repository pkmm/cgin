package service

import (
	"cgin/constant/miniprogram/menuaction"
	"cgin/model"
	"cgin/util"
	"github.com/pkg/errors"
	"sync"
)

type miniProgramService struct {
	mutex *sync.Mutex
}

var MiniProgramService = &miniProgramService{
	mutex: &sync.Mutex{},
}

var (
	menuAlreadyExist = errors.New("菜单项已经存在")
)

func (m *miniProgramService) DisposeMenu(desp, title, icon string, actionType menuaction.MiniProgramMenuActionType, actionValue string) (menu *model.Menu, err error) {

	if m.GetMenuByTitle(title) != nil {
		return nil, menuAlreadyExist
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	menu = &model.Menu{
		Desp:        desp,
		Title:       title,
		Icon:        icon,
		ActionType:  actionType,
		ActionValue: actionValue,
	}
	tx := db.Begin()
	if err = tx.Save(menu).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return menu, nil
}

func (m *miniProgramService) GetMenuByTitle(title string) *model.Menu {
	menu := &model.Menu{}
	if err := db.Model(&model.Menu{}).Where("title = ?", title).First(&menu).Error; err != nil {
		return nil
	}
	return menu
}

func (m *miniProgramService) GetAllMenus() []*model.Menu {
	var menus []*model.Menu
	if err := db.Model(&model.Menu{}).Find(&menus).Error; err != nil {
		return nil
	}
	return menus
}

func (m *miniProgramService) GetAllActiveMenus() []*model.Menu {
	var menus []*model.Menu
	if err := db.Model(&model.Menu{}).Where("disabled = 0").Find(&menus).Error; err != nil {
		return nil
	}
	return menus
}

func (m *miniProgramService) GetIndexConfig() *model.IndexConfig {
	indexConfig := &model.IndexConfig{}
	if err := db.Model(&model.IndexConfig{}).Where("disabled = 0").Order("id DESC").First(&indexConfig).Error; err != nil {
		return nil
	}
	return indexConfig
}

func (m *miniProgramService) SetIndexConfig(slogan, imageUrl, style string) *model.IndexConfig {
	indexConfig := &model.IndexConfig{
		Slogan:     slogan,
		ImageUrl:   imageUrl,
		ImageStyle: style,
		Disabled:   false,
	}

	if err := db.Save(indexConfig).Error; err != nil {
		return nil
	}
	return indexConfig
}

func (m *miniProgramService) UpdateNotification(id uint64, content string, startAt, endAt util.JSONTime) *model.Notification {
	if err := db.Model(&model.Notification{}).Where("id = ?", id).Update(map[string]interface{}{
		"content": content, "start_at": startAt, "end_at": endAt}).Error; err != nil {
		return nil
	}
	return m.GetNotificationById(id)
}

func (m *miniProgramService) GetNotificationById(id uint64) *model.Notification {
	notification := &model.Notification{}
	if err := db.Model(&model.Notification{}).Where("id = ?", id).First(&notification).Error; err != nil {
		return nil
	}
	return notification
}

func (m *miniProgramService) SaveNotification(content string, startAt, endAt util.JSONTime) *model.Notification {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	notification := &model.Notification{
		Content: content,
		StartAt: startAt,
		EndAt:   endAt,
	}
	tx := db.Begin()
	if err := tx.Save(notification).Error; err != nil {
		tx.Rollback()
		return nil
	}
	tx.Commit()
	return notification
}

func (m *miniProgramService) GetLatestNotification() *model.Notification {
	notification := &model.Notification{}
	if err := db.Model(&model.Notification{}).Order("id desc").
		Limit(1).First(&notification).Error; err != nil {
		return nil
	}
	return notification
}

func (m *miniProgramService) GetNotifications(limit uint64) []*model.Notification {
	var notifications []*model.Notification
	if err := db.Model(&model.Notification{}).
		Where("disabled = 0").
		Order("id desc").Limit(limit).Find(&notifications).Error; err != nil {
		return notifications
	}
	return notifications
}
