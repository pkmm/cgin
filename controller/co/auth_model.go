package co

import "cgin/constant/devicetype"

type AuthModel struct {
	Openid     string                `json:"openid" example:"aqe"`
	Token      string                `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjgsImV4cCI6MTU2MjQ2NTY0MSwiaXNzIjoiY2NsYSJ9.J3QXF3tZlvzyWvx8VG9EibUIr5mHK0xg3mjxY8LGhk8"`
	Sign       string                `json:"sign" example:"ewe"`
	DeviceType devicetype.DeviceType `json:"device_type" example:"1"`
}
