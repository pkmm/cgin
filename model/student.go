package model

type Student struct {
	Id       uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId   uint64 `json:"user_id" gorm:"index:idx_user_id;not null;"`
	Name     string `json:"name" gorm:"type:varchar(64);default:null;"`
	Number   string `json:"number" gorm:"default:null;type:varchar(64);"`
	Password string `json:"password" gorm:"default:null;type:varchar(64);"`
	IsSync   bool   `json:"is_sync" gorm:"default:1"`
	Model

	Scores     []*Score    `gorm:"ForeignKey:StudentId;AssociationForeignKey:Id" json:"scores"`
	SyncDetail *SyncDetail `json:"sync_detail" gorm:"ForeignKey:StudentId;AssociationForeignKey:Id"`
}
