package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"io/ioutil"
	"path"
	"encoding/json"
	"fmt"
	"net/url"
	"github.com/garyburd/redigo/redis"
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

type Config struct {
	MySQL MySQLConfig `json:"mysql"`
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
			"%s:%s@/%s?charset=utf8&parseTime=False&loc=%s",
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
	// mysql 连接池 配置
	db.DB().SetMaxIdleConns(5) // 空闲链接
	db.DB().SetMaxOpenConns(100)

	pool = newPool("127.0.0.1:6379")

}

func GetDB() *gorm.DB {
	db.DB().Ping()
	return db
}

// 创建redis的连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func GetRedis() redis.Conn {
	return pool.Get()
}
