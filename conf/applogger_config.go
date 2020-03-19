package conf

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
)

// 使用beego的log系统
var Logger = InitLogger()

func InitLogger() *logs.BeeLogger {
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

	wd, _ := os.Getwd()
	logFileStorageIn := path.Join(wd, "storage/logs/gin.log")
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
	fmt.Printf("log file storage in: %#v\n", logFileStorageIn)
	return logger
}
