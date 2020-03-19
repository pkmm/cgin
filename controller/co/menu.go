package co

import (
	"cgin/constant/miniprogram/menuaction"
)

// 小程序菜单的通信的数据结构
type Menu struct {
	Desp        string                    `json:"desp,omitempty" example:"这是一个菜单的描述"`
	Title       string                    `json:"title" example:"标题"`
	Icon        string                    `json:"icon,omitempty" example:"icon"`
	ActionType  menuaction.MenuActionType `json:"action_type" example:"2"`
	ActionValue string                    `json:"action_value" example:"action value"`
}

type Menus struct {
	Menus []Menu `json:"menus"`
}
