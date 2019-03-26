package service

import (
	"cgin/conf"
	"cgin/model"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/url"
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

var db *gorm.DB
var err error
var workPath string
var pool *redis.Pool
var dsn string

func init() {
	dsn = fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=%s",
		conf.AppConfig.String("mysql.username"),
		conf.AppConfig.String("mysql.password"),
		conf.AppConfig.String("mysql.database"),
		url.QueryEscape(conf.AppConfig.String("mysql.timezone")),
	)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Mysql 连接池 配置
	db.DB().SetMaxIdleConns(10) // 空闲链接
	db.DB().SetMaxOpenConns(100)

	if conf.AppConfig.String("appEnv") != "prod" {
		db.LogMode(true)
	}

	// == Redis == 配置
	// redis 连接池 设置
	pool = newPool("127.0.0.1:6379")

	ConnectDB()
}

func GetDB() *gorm.DB {
	return db
}

// 创建redis的连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func GetRedis() redis.Conn {
	return pool.Get()
}

func ConnectDB() {
	var err error

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(model.Models...).Error; err != nil {
		conf.AppLogger.Error("auto migrate tables failed, " + err.Error())
	}
}
