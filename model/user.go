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
