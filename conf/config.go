package conf

import (
	"cgin/util"
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"path"
)

// 使用beego的config包
var AppConfig config.Configer

// 使用beego的log系统
var Logger *logs.BeeLogger

const (
	// 运行环境的参数设置
	AppEnvironment = "appEnv"
	AppEnvProd     = "prod"
	AppEnvDev      = "dev"

	// mysql
	MysqlHost     = "mysql.host"
	MysqlPort     = "mysql.port"
	MysqlUser     = "mysql.username"
	MysqlPassword = "mysql.password"
	MysqlDatabase = "mysql.database"
	MysqlTimezone = "mysql.timezone"
)

func init() {
	// Note: Must init logger module first.
	Logger = initLogger()
	AppConfig = initAppConfig()
}

func initLogger() *logs.BeeLogger {
	// 设置beego logs
	logger := logs.NewLogger(1e5)
	logger.SetPrefix("[RETAIN]:")
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(20)
	logger.Async(1e3)

	type logConfig struct {
		Filename string   `json:"filename"`
		Level    int      `json:"level"`
		Daily    bool     `json:"daily"`
		Maxdays  int      `json:"maxdays"`
		Separate []string `json:"separate"`
	}

	wd := util.GetCurrentCodePath()
	logFileStorageIn := path.Join(wd, "..", "storage/logs/gin.log")
	lf := &logConfig{
		Filename: logFileStorageIn,
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
		logger.SetLogger(logs.AdapterMultiFile, string(byteOfLF))
	}

	// print config ini file path.
	logger.Info("log file storage in: %#v", logFileStorageIn)
	return logger
}

func initAppConfig() config.Configer {
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
