package model

import "cgin/model/modelInterface"

type Tieba struct {
	Id     uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId uint64 `json:"user_id" gorm:"index"`
	Bduss  string `json:"bduss" gorm:"default:null;type:varchar(255);"`
	Result string `json:"result" gorm:"default:null;"`
	User   *User  `json:"user" gorm:"ForeignKey:UserId;AssociationForeignKey:Id;"`
	Model
}

func (t *Tieba) GetList(info modelInterface.PageSizeInfo) (error, interface{}, int) {
	err, query, total := basicPagination(info, t)
	if err != nil {
		return err, nil, 0
	} else {
		var result []*Tieba
		err = query.Find(&result).Error
		return err, result, total
	}
}
