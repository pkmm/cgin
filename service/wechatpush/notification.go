package wechatpush

import "net/http"

// 策略
type Notify interface {
	Send(title, content string) (*http.Response, error)
}
