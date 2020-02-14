package service

import "cgin/model"

// 收集一些美好的有意义的句子，思考人生的意义是什么

type thinkingService struct {
	baseService
}

var ThinkingService = &thinkingService{}

func (s *thinkingService) GetOneById(id uint64) *model.Thinking {
	var result model.Thinking
	if err := db.Model(&model.Thinking{}).Where("id = ?", id).First(&result).Error; err != nil {
		return nil
	}
	return &result
}

func (s *thinkingService) GetOneByUserId(userId uint64) *model.Thinking {
	var result model.Thinking
	if err := db.Model(&model.Thinking{}).Where("user_id = ?", userId).First(&result).Error; err != nil {
		return nil
	}
	return &result
}

func (s *thinkingService) SaveOne(userId uint64, text string) *model.Thinking {
	t := model.Thinking{UserId: userId, Content: text}
	if err := db.Model(&model.Thinking{}).Save(t).Error; err != nil {
		return nil
	}
	return s.GetOneByUserId(userId)
}

func (s *thinkingService) RemoveOne(id uint64) bool {
	if err := db.Model(&model.Thinking{}).Where("id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func (s *thinkingService) GetList(page, size int) []*model.Thinking {
	var results []*model.Thinking
	if page < 1 || size <= 0 {
		panic("查询参数错误")
	}
	offset := (page - 1) * size
	db.Model(&model.Thinking{}).Where("is_deleted = ? or is_deleted is NULL", false).
		Limit(size).Offset(offset).Find(&results).Order("id asc")
	return results
}
