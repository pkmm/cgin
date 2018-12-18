package middleware

import (
	"bytes"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"net/http"
	"pkmm_gin/model"
	"pkmm_gin/utility"
	"strconv"
	"strings"
	"time"
)

// APi 认证
// 1. 使用账密码登陆，获取token
// 2. 使用token进行url参数签名
// 3. 服务端使用redis缓存token

type AuthData struct {
	UserId      int    `json:"uid"`
	AccessToken string `json:"tk"`
}

type AuthDataInfo struct {
	AuthData
	Ts    int64
	Nonce string
	Sign  string
}

func ApiAuth() gin.HandlerFunc {
	const TTL = 60 // 时间有效期是60s
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if strings.Index(url, "login") != -1 {
			c.Next()
			return
		}
		// 创建buffer 备份body
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		// 本次中间件的读取
		c.Request.Body = rdr1
		valid := true
		myRedis := model.GetRedis()
		defer myRedis.Close()
		var timestamp string
		var token string
		var nonce string
		// todo 下面还需要再想想
		if c.Request.Method == "POST" {
			var authDataInfo AuthDataInfo
			if err := c.Bind(&authDataInfo); err == nil {
				timestamp = strconv.FormatInt(authDataInfo.Ts, 10)
				token = authDataInfo.AccessToken
				nonce = authDataInfo.Nonce
				myAccessToken, er1 := redis.String(myRedis.Do("GET", authDataInfo.UserId))
				if er1 == redis.ErrNil || myAccessToken != token {
					valid = false
				}
				newSign := utility.Signature(map[string]string{
					"tk":    authDataInfo.AccessToken,
					"uid":   strconv.Itoa(authDataInfo.UserId),
					"ts":    strconv.FormatInt(authDataInfo.Ts, 10),
					"nonce": authDataInfo.Nonce,
				})
				if newSign != authDataInfo.Sign {
					valid = false
				}
			} else {
				valid = false
			}
		} else if c.Request.Method == "GET" {
			uid := c.DefaultQuery("uid", "")
			ts := c.DefaultQuery("ts", "")
			nonce = c.DefaultQuery("nonce", "")
			tk := c.DefaultQuery("tk", "")
			sign := c.DefaultQuery("sign", "")
			if uid == "" || ts == "" || nonce == "" || tk == "" || sign == "" {
				valid = false
			} else {
				timestamp = ts
				token = tk
				newSign := utility.Signature(map[string]string{
					"uid":   uid,
					"ts":    ts,
					"nonce": nonce,
					"tk":    tk,
				})
				if newSign != sign {
					valid = false
				}
			}
		}
		// 验证时间有效性
		if valid {
			now := time.Now().Unix()
			if ts, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
				if math.Abs(float64(now-ts)) > TTL {
					valid = false
				}
			} else {
				valid = false
			}
		}

		// 验证指定的时间内的唯一性
		// 维护nonce的集合，只要nonce在集合中存在就是不合法的，在维护一个有序的集合，然后定时
		// 删除过期的nonce
		if valid {
			tmpKey := token + nonce
			_, err := redis.Int64(myRedis.Do("ZSCORE", "api_collection_ttl", tmpKey))
			if err == redis.ErrNil {
				ts := time.Now().Add(time.Second * TTL).Unix()
				myRedis.Do("ZADD", "api_collection_ttl", ts, tmpKey)
			} else {
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

		// 重新设置body
		c.Request.Body = rdr2
		c.Next()
	}
}
