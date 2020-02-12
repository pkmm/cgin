package task

import (
	"cgin/conf"
	"cgin/service"
)

func backupMysql() {
	conf.Logger.Info("备份mysql数据库已经开始")
	service.BackupDBService.BackupMysqlService("")
}
