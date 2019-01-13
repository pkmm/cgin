package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"pkmm_gin/model"
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05"

var db *gorm.DB
var err error
var workPath string
var pool *redis.Pool
var dsn string

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	TimeZone string `json:"time_zone"`
}

type WeChatSmallProgramConfig struct {
	Secret string `json:"secret"`
	AppId  string `json:"app_id"`
}

// 新添加的配置项放在此处就能解析到
type Config struct {
	MySQL MySQLConfig `json:"mysql"`
	// 微信小程序的配置文件
	WeChatSmallProgram WeChatSmallProgramConfig `json:"wechat_small_program"`
}

var config Config

func init() {

	workPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	cfgPath := path.Join(workPath, "conf", "conf.json")
	cfgbuf, err := ioutil.ReadFile(cfgPath)
	if err = json.Unmarshal(cfgbuf, &config); err == nil {
		dsn = fmt.Sprintf(
			"%s:%s@/%s?charset=utf8&parseTime=True&loc=%s",
			config.MySQL.User,
			config.MySQL.Password,
			config.MySQL.Database,
			url.QueryEscape(config.MySQL.TimeZone),
		)
		db, err = gorm.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}

	// Mysql 连接池 配置
	db.DB().SetMaxIdleConns(10) // 空闲链接
	db.DB().SetMaxOpenConns(100)
	//f, _ := os.Create("gin.log")
	//db.SetLogger(log.New(f, "\r\n", 0))
	db.LogMode(true)

	// == Redis == 配置
	// redis 连接池 设置
	pool = newPool("127.0.0.1:6379")

	ConnectDB()

}

func GetDB() *gorm.DB {
	//db.DB().Ping()
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

func GetConfig() Config {
	return config
}

func ConnectDB() {
	var err error

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;").
		AutoMigrate(model.Models...).Error; err != nil {
		logs.Error("auto migrate tables failed, " + err.Error())
	}
}
