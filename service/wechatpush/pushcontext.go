package wechatpush

import "net/http"

/////////////////////
/// context上下文/////
/////////////////////
type PusherContext struct {
	n              Notify
	title, content string
}

func NewPusherContext(n Notify, title , content string) *PusherContext {
	return &PusherContext{n: n, title: title, content: content}
}

func (p *PusherContext) ChangeNotify(n Notify) *PusherContext {
	p.n = n
	return p
}

func (p *PusherContext) Push() (*http.Response, error) {
	return p.n.Send(p.title, p.content)
}
