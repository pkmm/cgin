package co

import (
	"cgin/constant/miniprogram/menuaction"
)

// 小程序菜单的通信的数据结构
type Menu struct {
	Desp        string                               `json:"desp,omitempty"`
	Title       string                               `json:"title"`
	Icon        string                               `json:"icon,omitempty"`
	ActionType  menuaction.MiniProgramMenuActionType `json:"action_type"`
	ActionValue string                               `json:"action_value"`
}
