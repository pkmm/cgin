package co

import "cgin/util"

type Notification struct {
	AuthCredit
	Id      uint64        `json:"id,omitempty"`
	Content string        `json:"content,omitempty" example:"lalala"`
	StartAt util.JSONTime `json:"start_at,omitempty" example:"1582044902000"`
	EndAt   util.JSONTime `json:"end_at,omitempty" example:"1676739295000"`
}
