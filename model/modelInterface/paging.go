package modelInterface

type (
	PaginatedModel interface {
		GetList(info PageSizeInfo) (err error, list interface{}, total int)
	}
	PageSizeInfo struct {
		Page, PageSize int
	}
)
