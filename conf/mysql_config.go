package conf

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/url"
)

var DB = InitMysql(GetMysqlConfigInfo(), IsDev())

func InitMysql(info *MysqlConfigInfo, debug bool) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=%s",
		info.User,
		info.Pwd,
		info.Database,
		info.Timezone,
	)
	DB, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Mysql 连接池 配置
	DB.DB().SetMaxIdleConns(10) // 空闲链接
	DB.DB().SetMaxOpenConns(100)

	if debug {
		DB.LogMode(true)
	}
	Logger.Info("mysql连接成功")
	return DB
}

type MysqlConfigInfo struct {
	Host, Port, User, Pwd, Database, Timezone string
}

func GetMysqlConfigInfo() *MysqlConfigInfo {
	m := &MysqlConfigInfo{
		Host:     AppConfig.String(mysqlHost),
		Port:     AppConfig.String(mysqlPort),
		User:     AppConfig.String(mysqlUser),
		Pwd:      AppConfig.String(mysqlPassword),
		Database: AppConfig.String(mysqlDatabase),
		Timezone: url.QueryEscape(AppConfig.String(mysqlTimezone)),
	}
	return m
}
