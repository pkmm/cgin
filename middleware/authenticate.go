package middleware

import (
	"cgin/errno"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

// 使用JWT进行认证

const UID = "uid"
const Token = "token"

type authData struct {
	Token string `json:"token"`
}

func Auth(c *gin.Context) {
	// 创建buffer 备份body
	//buf, _ := ioutil.ReadAll(c.Request.Body)
	//rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	//rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
	//
	//c.Request.Body = rdr1
	//data := new(authData)
	//if c.Request.Method == http.MethodPost {
	//	if err := c.ShouldBindJSON(&data);err != nil {
	//		panic(errno.NormalException.AppendErrorMsg("获取token失败"))
	//	}
	//} else if c.Request.Method == http.MethodGet {
	//	data.Token = c.DefaultQuery(Token, "")
	//	if data.Token == "" {
	//		panic(errno.TokenNotValid)
	//	}
	//}
	//c.Request.Body = rdr2
	token := c.Request.Header.Get("Authorization")
	if 0 == len(token) {
		panic(errno.TokenNotValid)
	}
	claims, err := service.JWTSrv.GetAuthClaims(token)
	if err != nil {
		panic(errno.TokenNotValid)
	}
	c.Set(UID, claims.Uid)
	c.Next()
}
