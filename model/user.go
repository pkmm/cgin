package model

type User struct {
	Id     uint64 `json:"id" gorm:"primary_key"`
	OpenId string `json:"open_id" gorm:"unique;size:128;default:null"`

	// 指针类型 保证关联数据不存在的时候能正常显示为null
	Student *Student `json:"student" gorm:"ForeignKey:UserId;AssociationForeignKey:Id"`
	Model
}
