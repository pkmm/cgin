package resposne

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS = 0
	ERROR   = 76
)

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{data, code, msg})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, "OK", nil, c)
}

func OkWithMsg(msg string, c *gin.Context) {
	Result(SUCCESS, msg, nil, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, "OK", data, c)
}

func OkWithDetail(msg string, data interface{}, c *gin.Context) {
	Result(SUCCESS, msg, data, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, "ERROR", nil, c)
}

func FailWithMsg(msg string, c *gin.Context) {
	Result(ERROR, msg, nil, c)
}

func FailWithDetail(msg string, data interface{}, c *gin.Context) {
	Result(ERROR, msg, data, c)
}
