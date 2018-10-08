package middleware

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkmm_gin/model"
	"pkmm_gin/utility"
	"strconv"
	"strings"
	"time"
)

// APi 认证
// 1. 使用账密码登陆，获取code
// 2. 使用code进行url参数签名
// 3. 服务端使用redis缓存code

func ApiAuth() gin.HandlerFunc {
	const COLLECTION = "api_collection"
	const TTL = 60 // 时间有效期是60s
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if strings.Index(url, "student/login") != -1 {
			c.Next()
			return
		}
		signature := c.DefaultQuery("signature", "")
		timestamp := c.DefaultQuery("timestamp", "")
		nonce := c.DefaultQuery("nonce", "")
		code := c.DefaultQuery("code", "")
		valid := true
		// 验证参数
		if signature == "" || timestamp == "" || code == "" || nonce == "" {
			valid = false
		}
		// 验证时间有效性
		if valid {
			now := time.Now().Unix()
			if ts, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
				if now < ts || now-ts > TTL {
					valid = false
				}
			} else {
				valid = false
			}
		}

		myRedis := model.GetRedis()
		defer myRedis.Close()

		// 验证指定的时间内的唯一性
		// 维护nonce的集合，只要nonce在集合中存在就是不合法的，在维护一个有序的集合，然后定时
		// 删除过期的nonce
		if valid {
			exists, _ := redis.Bool(myRedis.Do("SISMEMBER", COLLECTION, nonce))
			if exists {
				valid = false
			} else {
				myRedis.Do("SADD", COLLECTION, nonce)
				ts := time.Now().Add(time.Second * TTL).Unix()
				myRedis.Do("ZADD", "api_collection_ttl", ts, nonce)
			}
		}

		key, _ := redis.String(myRedis.Do("GET", code))

		// 验证签名值
		if valid {
			sign := utility.Signature(map[string]string{
				"timestamp":    timestamp,
				"code":         code,
				"security_key": key,
				"nonce":        nonce,
			})
			if sign != signature {
				valid = false
			}
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid request.",
			})
			c.Abort() // Stop execute next middleware.
			return
		}

		c.Next()
	}
}
