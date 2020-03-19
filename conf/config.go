package conf

import (
	"cgin/util"
	"github.com/astaxie/beego/config"
	"path"
)

// 使用beego的config包
var AppConfig = InitAppConfig()

func InitAppConfig() config.Configer {
	wd := util.GetCurrentCodePath()
	confPath := path.Join(wd, ".env")
	myAppConfig, err := config.NewConfig("ini", confPath)
	if err != nil {
		panic(err)
	}
	if Logger != nil {
		Logger.Info("conf path: %#v", confPath)
	}
	return myAppConfig
}
