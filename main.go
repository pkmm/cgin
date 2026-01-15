package main

import (
	"cgin/core"
	"cgin/global"
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
	global.GLog = core.Zap()
}
