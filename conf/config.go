package conf

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

// 使用beego的config包
var AppConfig config.Configer
var err error

func init() {
	AppConfig, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println(err)
		panic("Error of load app.conf")
	}
}
