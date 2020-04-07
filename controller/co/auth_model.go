package co

import "cgin/constant/devicetype"

// deviceType是2的时候使用openid 和sign字段
// 是3的时候使用username 和 password字段
type AuthModel struct {
	Openid     string                `json:"openid" example:"openid_xxsd"`
	Sign       string                `json:"sign" example:"67807AFF5A99880726B74D03F5A8F78C"`
	Username   string                `json:"username" example:"cc"`
	Password   string                `json:"password" example:"x"`
	DeviceType devicetype.DeviceType `json:"device_type" example:"2"`
}
