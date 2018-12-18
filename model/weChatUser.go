package model

import "time"

// 微信信息
type WeChatUser struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	AvatarUrl string    `json:"avatar" gorm:"column:avatar"`
	Country   string    `json:"country"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Nickname  string    `json:"nickname" gorm: "column:nickname"`
	Language  string    `json:"language"`
	Gender    int       `json:"gender,string"`
	UnionId   string    `json:"union_id"`
	OpenId    string    `json:"open_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 使用code换取sessionKey
type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
	Errcode    string `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

// 微信小程序登陆的请求参数
type WxLoginRequest struct {
	Iv   string `json:"iv" binding:"required"`
	Data string `json:"encrypted_data" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func (WeChatUser) TableName() string {
	return "wechat_users"
}

// update or create record.
func CreateWeChatUser(u WeChatUser) WeChatUser {
	db.Where(WeChatUser{OpenId: u.OpenId}).Assign(u).FirstOrCreate(&u)
	return u
}

func GetWeChatUserByOpenId(openId string) WeChatUser {
	var wechatUser WeChatUser
	db.Where(WeChatUser{OpenId: openId}).First(&wechatUser)
	return wechatUser
}
