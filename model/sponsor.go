package model

// 赞助者
type Sponsor struct {
	Id         uint64 `json:"id" gorm:"primary_key;auto_increment;"`
	UserId     uint64 `json:"user_id" gorm:"index"`
	Money      uint   `json:"money" gorm:"default:0;"` // 捐赠的钱，单位是分
	PayChannel int    `json:"pay_channel" gorm:"default:null;"` // 支付的渠道
	User       *User  `json:"user" gorm:"ForeignKey:UserId;AssociationForeignKey:Id;"`
	Model
}

const (
	PAY_CHANNEL_UNKNOW = iota
	PAY_CHANNEL_WECHAT
	PAY_CHANNEL_ALIPAY
)