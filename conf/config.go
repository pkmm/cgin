package conf

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

// 使用beego的config包
var AppConfig config.Configer
var err error

// 使用beego的log系统
var AppLogger *logs.BeeLogger

// 运行环境的参数设置
const (
	AppEnvironment = "appEnv"
	EnvProd = "prod"
	EnvDev = "dev"
)

func init() {
	AppConfig, err = config.NewConfig("ini", "conf/.env")
	if err != nil {
		fmt.Println(err)
		panic("Error of load .env")
	}

	// 设置beego logs
	AppLogger = logs.NewLogger(1e5)
	AppLogger.SetPrefix("[USE GIN]:")
	AppLogger.EnableFuncCallDepth(true)
	AppLogger.SetLogFuncCallDepth(10)
	AppLogger.Async(1e3)

	AppLogger.SetLogger(logs.AdapterMultiFile, `{
	"filename":"storage/logs/cgin.log",
	"level":7,
	"daily":true,
	"maxdays":2,
	"separate": ["error", "info", "debug"]}`)
}
