package conf

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool = newPool("127.0.0.1:6379")

func GetRedis() redis.Conn {
	return pool.Get()
}

// 创建redis的连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}
