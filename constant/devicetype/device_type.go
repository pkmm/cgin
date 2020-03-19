package devicetype

// 请求的客户端设备
type DeviceType int

const (
	Unknow DeviceType = iota + 1
	MiniProgram
	WebBrowser
)
