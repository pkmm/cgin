package model

import "cgin/conf"

type User struct {
	Id     uint64 `json:"id" gorm:"primary_key"`
	OpenId string `json:"open_id" gorm:"unique;size:128;default:null"`
	Roles  []Role `json:"roles" gorm:"many2many:user_roles;"`

	// 指针类型 保证关联数据不存在的时候能正常显示为null
	Student *Student `json:"student" gorm:"ForeignKey:UserId;AssociationForeignKey:Id"`
	Model
}

func GetUserById(userId uint64) (error, *User) {
	var user User
	err := conf.DB.Find(&user, User{Id: userId}).Error
	return err, &user
}

func (u *User) GetRoles() (err error, roles *[]Role) {
	var r []Role
	err = conf.DB.Model(u).Related(&r, "Roles").Error
	return err, &r
}

func (u *User) GetStudent() (error, *Student) {
	var stu Student
	err := conf.DB.Model(u).Related(&stu).Error
	return err, &stu
}
func IsAdmin(userId uint64) bool {
	var t int = 0
	conf.DB.Table("user_roles").
		Where("user_id = ? AND role_id = ?", userId, AdminRoleId()).Count(&t)
	return t > 0
}