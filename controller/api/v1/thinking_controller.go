package v1

import (
	"cgin/controller/contextHelper"
	"cgin/errno"
	"cgin/model"
	"github.com/gin-gonic/gin"
)

type thinkController struct{}

var ThinkingController = new(thinkController)

// @Summary 值得深思的句子
// @Tags 思考
// @Accept json
// @Produce json
// @Router /thinking [get]
// @Success 200 object service.Response
// @Param paging query co.PageLimitOffset true "page size"
func (t *thinkController) Index(ctx *gin.Context) {
	helper := contextHelper.New(ctx)
	page := helper.GetInt("page")
	size := helper.GetInt("size")
	err, results, total := new(model.Thinking).GetList(page, size)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	helper.Response(gin.H{
		"sentence": results,
		"total":    total,
	})
}
