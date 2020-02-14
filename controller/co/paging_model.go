package co

type PageLimitOffset struct {
	Page int `json:"page" example:"1"`
	Size int `json:"size" example:"100"`
}
