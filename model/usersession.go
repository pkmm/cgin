package model

type UserSession struct {
	ID          uint64 `json:"id" gorm:"primary_key"`
	UserId      uint64 `json:"user_id" gorm:"unique"`
	AccessToken string `json:"access_token" gorm:"size:128"`
	Active      bool   `json:"active" gorm:"default:1"`
	Model
}
