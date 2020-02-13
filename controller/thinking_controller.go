package controller

import (
	"cgin/controller/context_helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

type thinkController struct{}

var ThinkingController = new(thinkController)

func (t *thinkController) GetOne(ctx *gin.Context) {

}

func (t *thinkController) GetList(ctx *gin.Context) {
	helper := context_helper.New(ctx)
	page := helper.GetInt("page")
	size := helper.GetInt("size")
	fmt.Println("===>", page, size)
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//results := service.ThinkingService.GetList(page, size)
	helper.Response(fmt.Sprintf("page: %d, size: %d", page, size))
}
