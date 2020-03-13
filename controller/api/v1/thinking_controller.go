package v1

import (
	"cgin/controller/context_helper"
	"cgin/service"
	"github.com/gin-gonic/gin"
)

type thinkController struct{}

var ThinkingController = new(thinkController)

func (t *thinkController) GetOne(ctx *gin.Context) {

}

// @Summary 值得深思的句子
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Router /thinking/list [post]
// @Success 200 {object} service.Response
// @Param paging body co.PageLimitOffset true "page size"
func (t *thinkController) GetList(ctx *gin.Context) {
	helper := context_helper.New(ctx)
	page := helper.GetInt("page")
	size := helper.GetInt("size")
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	results := service.ThinkingService.GetList(page, size)
	helper.Response(results)
}
