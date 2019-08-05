package service

import "cgin/model"

type tiebaService struct {
	baseService
}

var TiebaService = &tiebaService{}

// TODO: 优化 分段取， 取某一列数据
func (t *tiebaService) GetAll() []*model.Tieba {
	var result []*model.Tieba
	if err := db.Model(&model.Tieba{}).Find(&result).Error; err != nil {
		return nil
	}
	return result
}

func (t *tiebaService) UpdateResultField(userId uint64, result string) *model.Tieba {
	var updatedModel model.Tieba
	if err := db.Model(&model.Tieba{}).Where("user_id = ?", userId).UpdateColumn("result", result).Error; err != nil {
		return nil
	}
	return &updatedModel
}
