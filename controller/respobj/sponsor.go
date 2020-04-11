package respobj

import "cgin/util"

type Sponsor struct {
	Id        uint64        `json:"id"`
	Money     uint          `json:"money"`
	OpenId    string        `json:"open_id"`
	Username  string        `json:"username"`
	CreatedAt util.JSONTime `json:"created_at"`
}
