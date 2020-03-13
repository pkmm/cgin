package model_interface

type PageSizeInfo struct {
	Page, PageSize int
}

type Paging interface {
	GetList(info PageSizeInfo) (err error, list interface{}, total int)
}