package service

import (
	"cgin/conf"
	"cgin/util"
	"io/ioutil"
	"os/exec"
	"path"
	"time"
)

type backupService struct {
	baseService
}

var BackupDBService baseService

func (b *baseService) BackupMysqlService(savePath string) {
	info := conf.GetMysqlConfigInfo()
	cmd := exec.Command("mysqldump", "-h"+info.Host, "-P"+info.Port, "-u"+info.User, "-p"+info.Pwd, info.Database)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic("备份mysql失败" + err.Error())
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	nowTime := time.Now().Format("2006_01_02_15_04_05")
	if savePath == "" {
		savePath = path.Join(util.GetCurrentCodePath(), "..", "storage", "backup", "mysql_db_backup_"+nowTime+".sql")
	}
	err = ioutil.WriteFile(savePath, bytes, 7744)
	if err != nil {
		panic(err)
	}
}
