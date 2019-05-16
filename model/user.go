package model

type User struct {
	Id       uint64 `json:"id" gorm:"primary_key"`
	OpenId   string `json:"open_id" gorm:"unique;size:128;default:null"`
	Model

	Student Student `json:"student" gorm:"ForeignKey:UserId:AssociationForeignKey:Id"`
}
