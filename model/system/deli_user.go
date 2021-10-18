package system

import "cgin/global"

type DeliUser struct {
	global.GModel
	Username string `json:"username" gorm:"type:varchar(20);comment:用户名拼音缩写;uniqueIndex:idx_n"`
	Token    string `json:"token" gorm:"comment:用户得力的token"`
	Cancel   int    `json:"cancel" gorm:"default:0;comment:不使用自动签到"`
	Uid      string `json:"uid" gorm:"comment:wxpusher UID"`
}

type DeliLoginData struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
