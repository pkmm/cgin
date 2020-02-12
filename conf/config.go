package conf

import (
	"cgin/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"path"
)

// 使用beego的config包
var AppConfig config.Configer
var err error

// 使用beego的log系统
var Logger *logs.BeeLogger

// 运行环境的参数设置
const (
	AppEnvironment = "appEnv"
	AppEnvProd     = "prod"
	AppEnvDev      = "dev"
)

func init() {
	wd := util.GetCurrentCodePath()
	fmt.Printf("conf path: %#v\n", wd)
	AppConfig, err = config.NewConfig("ini", path.Join(wd, ".env"))
	if err != nil {
		fmt.Println(err)
		panic("Error of load .env")
	}

	// 设置beego logs
	Logger = logs.NewLogger(1e5)
	Logger.SetPrefix("[USE GIN]:")
	Logger.EnableFuncCallDepth(true)
	Logger.SetLogFuncCallDepth(10)
	Logger.Async(1e3)

	type logConfig struct {
		Filename string   `json:"filename"`
		Level    int      `json:"level"`
		Daily    bool     `json:"daily"`
		Maxdays  int      `json:"maxdays"`
		Separate []string `json:"separate"`
	}

	lf := &logConfig{
		Filename: path.Join(wd, "..", "storage/logs/gin.log"),
		Level:    7,
		Daily:    true,
		Maxdays:  2,
		Separate: []string{"error", "info", "debug"},
	}
	//`{
	//"filename":"storage/logs/cgin.log",
	//"level":7,
	//"daily":true,
	//"maxdays":2,
	//"separate": ["error", "info", "debug"]}`

	if byteOfLF, err := json.Marshal(lf); err != nil {
		panic(err)
	} else {
		Logger.SetLogger(logs.AdapterMultiFile, string(byteOfLF))
	}
}
