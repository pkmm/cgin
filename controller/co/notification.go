package co

import "cgin/util"

type Notification struct {
	AuthCredit
	Id      uint64        `json:"id,omitempty" example:"18"`
	Content string        `json:"content,omitempty" example:"lalala"`
	StartAt util.JSONTime `json:"start_at,omitempty" example:"2019-01-23 23:23:34"`
	EndAt   util.JSONTime `json:"end_at,omitempty" example:"2029-01-23 23:12:30"`
}
