package model

import "cgin/conf"

type Resource struct {
	ID          int    `json:"id" gorm:"primary_key;"`
	Name        string `json:"name" gorm:"type:char(32);"`
	Description string `json:"description" gorm:"type:char(64);"`
	Identity    string `json:"identity" gorm:"type:char(32);unique;"`

	Permissions []Permission `json:"permissions" gorm:"ForeignKey:ResourceId;AssociationForeignKey:ID;"`
	Model
}

func (r *Resource) GetPermissions() (error, *[]Permission) {
	var t []Permission
	err := conf.DB.Model(r).Related(&t).Error
	return err, &t
}

func GetResourceByIdentity(identity string) (error, *Resource) {
	var t Resource
	err := conf.DB.Where(Resource{Identity: identity}).Find(&t).Error
	return err, &t
}
