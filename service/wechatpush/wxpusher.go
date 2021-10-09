package wechatpush

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
)

//////////////////////
// 使用wx pusher 推送//
//////////////////////
type wxPushContentType int

const (
	PlainText wxPushContentType = iota + 1
	Html
	Markdown
)

type pushBear struct {
	uids        []string
	contentType wxPushContentType
}

func NewPushBear(uids []string, contentType wxPushContentType) *pushBear {
	return &pushBear{uids: uids, contentType: contentType}
}

func (p *pushBear) Send(title, desc string) (*http.Response, error) {
	appToken := ""
	data := struct {
		AppToken    string            `json:"appToken"`
		Content     string            `json:"content"`
		ContentType wxPushContentType `json:"contentType"`
		TopicIds    []int             `json:"topicIds,omitempty"`
		Uids        []string          `json:"uids,omitempty"`
		Url         string            `json:"url,omitempty"`
	}{
		AppToken:    appToken,
		Content:     title + ", " + desc,
		ContentType: p.contentType,
		Uids:        p.uids,
	}
	if resp, _, errs := gorequest.New().
		Post("http://wxpusher.zjiecode.com/api/send/message").
		Send(data).
		End(); len(errs) != 0 {
		return nil, errs[0]
	} else {
		return resp, nil
	}
}
