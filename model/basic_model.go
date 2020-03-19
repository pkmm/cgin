package model

import (
	"cgin/conf"
	"cgin/model/modelInterface"
	"cgin/util"
	"github.com/jinzhu/gorm"
)
/**
 需要注意的是gorm对于默认空值在查询的时候是忽略的，例如查询disabled=false,如果使用结构体
	是无效的，因为怎么默认的空值行为是一致的
 */
type Model struct {
	CreatedAt util.JSONTime `json:"created_at" gorm:"type:timestamp;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt util.JSONTime `json:"updated_at" gorm:"type:timestamp;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//DeletedAt *time.Time `json:"-" gorm:"index:idx_deleted_at"`
}

func basicPagination(info modelInterface.PageSizeInfo, model modelInterface.PaginatedModel) (err error, query *gorm.DB, total int) {
	// default query ten record of per page.
	if info.Page <=0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	err = conf.DB.Model(model).Count(&total).Error
	query = conf.DB.Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).Order("id DESC")
	return err, query, total
}