package service

import (
	"cgin/errno"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, msg := errno.DecodeErr(err)
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func SendResponseWithInvalidParameters(c *gin.Context, data interface{}) {
	SendResponse(c, errno.InvalidParameters, data)
}

func SendResponseSuccess(c *gin.Context, data interface{}) {
	SendResponse(c, errno.Success, data)
}

func SendResponseWithInternalError(c *gin.Context) {
	SendResponse(c, errno.InternalServerError, nil)
}