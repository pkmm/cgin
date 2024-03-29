package main

import (
	"cgin/core"
	"cgin/global"
	"cgin/initialize"
	"cgin/schedule"
	"cgin/service/workerpool"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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
	global.GLog = core.Zap()
	global.DB = initialize.Gorm()

	// 初始化数据库
	if global.DB != nil {
		initialize.MysqlTables(global.DB) // 注册所有的表
		db, _ := global.DB.DB()
		defer db.Close()
	} else {
		_ = fmt.Errorf("初始化数据库失败!")
	}

	// 初始化worker pool
	pool, err := workerpool.NewPool(20, time.Second*10)
	if err != nil {
		panic("initialize worker pool failed.")
	}
	defer pool.Close()
	global.WorkerPool = pool

	// 初始化任务调度
	schedule.SC = schedule.NewSchedule()
	schedule.SC.StartJobs()
	defer schedule.SC.Stop()

	gin.SetMode(gin.ReleaseMode)
	core.RunWindowsServer()
}
