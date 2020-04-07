package model

import "cgin/conf"

const (
	RoleAdmin = iota + 1
	RoleUser
)

type Role struct {
	ID          int          `json:"id" gorm:"primary_key;"`
	Name        string       `json:"name" gorm:"type:char(32);default:null;"`
	Description string       `json:"description" gorm:"type:char(64);default:null;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Model
}

// 角色的所有权限
func (r *Role) GetPermissions() (error, *[]Permission) {
	var t []Permission
	err := conf.DB.Model(r).Related(&t, "Permissions").Error
	return err, &t
}

//
func HasPermission(permissionId int, roleId int) bool {
	var count int = 0
	conf.DB.Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		Count(&count)
	return count != 0
}

func AdminRoleId() int {
	var r Role
	conf.DB.Where(Role{ID: RoleAdmin}).Find(&r)
	return r.ID
}
