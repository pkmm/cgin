package main

import (
	"cgin/core"
	"cgin/global"
	"cgin/initialize"
	"fmt"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title cgin server api
// @version 1.0
// @description 实验性

// @host localhost:8654
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	// 初始化需要的配置信息
	global.G_VP = core.Viper()
	global.G_LOG = core.Zap()
	global.G_DB = initialize.Gorm()

	if global.G_DB != nil {
		initialize.MysqlTables(global.G_DB) // 注册所有的表
		db, _ := global.G_DB.DB()
		defer db.Close()
	} else {
		_ = fmt.Errorf("初始化数据库失败!")
	}

	core.RunWindowsServer()
}
