package system

import "cgin/global"

type BaiduUser struct {
	global.GModel
	Name  string `json:"name" gorm:"comment:用户名"`
	Bduss string `json:"bduss" gorm:"comment:百度token"`
}
