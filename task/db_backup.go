package task

import (
	"cgin/conf"
	"cgin/service"
	"time"
)

func backupMysql() {
	conf.Logger.Info("备份mysql数据库已经开始")
	st := time.Now()
	service.BackupDBService.BackupMysqlService("")
	conf.Logger.Info("mysql数据库备份完成: 用时：%s", time.Since(st).String())
}
