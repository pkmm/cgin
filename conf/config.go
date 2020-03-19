package conf

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"os"
	"path"
)

// 使用beego的config包
var AppConfig = InitAppConfig()

func InitAppConfig() config.Configer {
	wd, _ := os.Getwd()
	confPath := path.Join(wd, "conf/.env")
	myAppConfig, err := config.NewConfig("ini", confPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("conf path: %#v\n", confPath)
	return myAppConfig
}
