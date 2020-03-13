package service

import (
	"cgin/model/model_interface"
	"github.com/jinzhu/gorm"
)

type pagingService struct {}

var PagingService pagingService

func (p *pagingService) GetList(model model_interface.Paging, info model_interface.PageSizeInfo) (err error, query *gorm.DB, total int) {
	// default query ten record of per page.
	if info.Page <=0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	err = db.Where(model).Count(&total).Error
	query = db.Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).Order("id DESC")
	return err, query, total
}
