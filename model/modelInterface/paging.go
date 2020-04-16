package modelInterface

type (
	PaginatedModel interface {
		GetList(page, size int) (err error, list interface{}, total int)
	}
)
