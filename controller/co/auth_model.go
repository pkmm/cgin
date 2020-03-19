package co

import "cgin/constant/devicetype"

type AuthModel struct {
	Openid     string                `json:"openid" example:"openid_xxsd"`
	Sign       string                `json:"sign" example:"67807AFF5A99880726B74D03F5A8F78C"`
	DeviceType devicetype.DeviceType `json:"device_type" example:"2"`
}
