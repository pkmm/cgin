package co

import "cgin/util"

type Notification struct {
	Id      uint64        `json:"id,omitempty"`
	Content string        `json:"content,omitempty"`
	StartAt util.JSONTime `json:"start_at,omitempty"`
	EndAt   util.JSONTime `json:"end_at,omitempty"`
}
