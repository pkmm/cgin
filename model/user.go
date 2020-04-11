package model

import "cgin/conf"

type User struct {
	Id       uint64 `json:"id" gorm:"primary_key"`
	OpenId   string `json:"open_id" gorm:"unique;size:128;"` // openid 小程序，web端由服务端生成
	RoleId   int    `json:"role_id" gorm:"index:idx_role_id;default:2;"`
	Username string `json:"username" gorm:"size:64;default:null;index:ind_username;"`
	Password string `json:"password" gorm:"size:64;not null;"`

	Role *Role `json:"role" gorm:"ForeignKey:RoleId;AssociationForeignKey:Id;"`

	// 指针类型 保证关联数据不存在的时候能正常显示为null
	Student *Student `json:"student" gorm:"ForeignKey:UserId;AssociationForeignKey:Id"`
	Model
}

func (u *User) Create() *User {
	conf.DB.Create(u)
	return u
}

func GetUserByOpenId(openid string) *User {
	var user User
	if err := conf.DB.
		Preload("Student").
		Preload("Role").
		Where("open_id = ?", openid).
		First(&user).Error; err != nil {
		return nil
	}
	return &user
}

func GetUserByUsername(username string) *User {
	var user User
	if err := conf.DB.
		Preload("Student").
		Preload("Role").
		First(&user, "username = ?", username).Error; err != nil {
		return nil
	}
	return &user
}

func GetUserById(userId uint64) (error, *User) {
	var user User
	err := conf.DB.
		Preload("Student").
		Preload("Role").
		First(&user, User{Id: userId}).
		Error
	return err, &user
}

func (u *User) GetRole() (err error, roles *Role) {
	var r Role
	err = conf.DB.Model(u).Related(&r).Error
	return err, &r
}

func (u *User) GetStudent() (error, *Student) {
	var stu Student
	err := conf.DB.Model(u).Related(&stu).Error
	return err, &stu
}

func IsAdmin(userId uint64) bool {
	var t int = 0
	conf.DB.Table("users").Where("id = ?", userId).Select("role_id").Scan(&t)
	return t == RoleAdmin
}
