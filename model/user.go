package model

type User struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	UnionId  string `json:"union_id" gorm:"unique;size:128"`
	OpenId   string `json:"open_id" gorm:"unique;size:128"`
	Nickname string `json:"nickname" gorm: "column:nickname;size:64"` // 微信的昵称

	Num         string `json:"num" gorm:"default:null;size:64"` // tags "json:,string" 表示别的类型也能解析到
	Pwd         string `json:"pwd" grom:"default:null;size:64"` // 教务系统的密码
	StudentName string `json:"name" gorm:"size:64"`             // 正方教务系统的学生姓名
	CanSync     int    `json:"can_sync" gorm:"index;default:1"` // 是否能同步学生的成绩标记

	Model

	Scores []Score `gorm:"ForeignKey:UserId;AssociationForeignKey:Id" json:"scores"`
}

//// 登陆session认证
//// session
//type Session struct {
//	AccessToken string `json:"access_token"`
//	UserId      int64  `json:"user_id"`
//}

//func (User) TableName() string {
//	return "users"
//}

//
//func CreateUser(u User) User {
//	db.Create(&u)
//	return u
//}
//
//func GetUserById(id int64) User {
//	var u User
//	db.Where(User{Id: id}).First(&u)
//	return u
//}
//
//func GetUserByOpenId(openId string) (User, bool) {
//	var u User
//	db.Table("users").Joins("JOIN wechat_users AS WU ON WU.user_id = users.id").
//		Where("WU.open_id = ?", openId).Preload("Student").First(&u)
//
//	return u, u.Id > 0
//}
