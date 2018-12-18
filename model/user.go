package model

import "time"

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Salt      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Student Student `json:"student" gorm:"ForeignKey:UserId;AssociationForeignKey:Id;save_associations:false"`
}

// 登陆session认证
// session
type Session struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
}

func (User) TableName() string {
	return "users"
}

func CreateUser(u User) User {
	db.Create(&u)
	return u
}

func GetUserById(id int64) User {
	var u User
	db.Where(User{Id: id}).First(&u)
	return u
}

func GetUserByOpenId(openId string) (User, bool) {
	var u User
	db.Table("users").Joins("JOIN wechat_users AS WU ON WU.user_id = users.id").
		Where("WU.open_id = ?", openId).Preload("Student").First(&u)

	return u, u.Id > 0
}
