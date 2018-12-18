package myConfig

import "github.com/astaxie/beego/config"

// 使用beego的config包
var AppConfig config.Configer
var err error

func init() {
	AppConfig, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic("Error of init myConfig model.")
	}
}

