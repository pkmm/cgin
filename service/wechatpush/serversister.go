package wechatpush

import (
	"net/http"
	"net/url"
)

///////////////////////////
// 使用 server酱 推送///////
///////////////////////////
type pusherSister struct {
}

func NewPusherSister() *pusherSister {
	return &pusherSister{}
}

func (p *pusherSister) Send(title, desc string) (*http.Response, error) {
	key := ""
	client := &http.Client{}
	if resp, err := client.PostForm("https://sc.ftqq.com/"+key+".send", url.Values{
		"text": {title},
		"desp": {desc},
	}); err != nil {
		return nil, err
	} else {
		return resp, nil
	}
}
