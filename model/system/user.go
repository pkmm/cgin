package system

import "cgin/global"

type DeliUser struct {
	global.GModel
	Username string `json:"username" gorm:"comment:用户名拼音缩写"`
	Token    string `json:"token" gorm:"comment:用户得力的token"`
}
