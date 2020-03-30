package model

import "cgin/conf"

const (
	EffectAllow = iota + 1
	EffectOwner
)

type Permission struct {
	ID          int    `json:"id" gorm:"primary_key"`
	ResourceId  int    `json:"resource_id" gorm:"index;"`
	Name        string `json:"name" gorm:"type:char(32);"`
	Description string `json:"description" gorm:"type:char(64)"`
	Method      string `json:"method"` // http method
	Effect      int    `json:"effect"` // 作用于自己 还是 全部
	Model
}

func GetPermissionByResourceIdentityAndMethod(resourceIdentity string, method string) (error, *Permission) {
	var p Permission
	err := conf.DB.Table("resources").
		Joins("LEFT JOIN permissions ON resources.id = permissions.resource_id AND permissions.method = ?", method).
		Where("resources.identity = ?", resourceIdentity).
		Select("permissions.*").
		Find(&p).Error
	return err, &p
}