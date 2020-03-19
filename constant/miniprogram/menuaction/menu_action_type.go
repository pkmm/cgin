package menuaction

// 点击小程序菜单的行为
type MenuActionType int

const (
	Unknow     MenuActionType = iota // 无操作
	AlertModal                       // 弹出模态框
	GotoPage                         // 跳转到指定的页面
)
