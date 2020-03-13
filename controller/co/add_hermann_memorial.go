package co

import "cgin/util"

type AddHermannMemorial struct {
	Unit      uint          `json:"unit" example:"2"`
	TotalUnit uint          `json:"total_unit" example:"39"`
	StartAt   util.JSONTime `json:"start_at,string" example:"2019-01-23 23:12:30"`
}
