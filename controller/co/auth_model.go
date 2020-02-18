package co

import "cgin/constant/devicetype"

type AuthModel struct {
	Openid     string                `json:"openid" example:"openid_xxsd"`
	Sign       string                `json:"sign" example:"559EB671B3E536C8B1705EA9BD90FCFE"`
	DeviceType devicetype.DeviceType `json:"device_type" example:"2"`
}
