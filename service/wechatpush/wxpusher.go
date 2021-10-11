package wechatpush

import (
	"cgin/global"
	"fmt"
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
	TopicIds    []int
	contentType wxPushContentType
}

func NewPushBear(topicIds []int, contentType wxPushContentType) *pushBear {
	return &pushBear{TopicIds: topicIds, contentType: contentType}
}

func (p *pushBear) Send(title, desc, html string) (*http.Response, error) {
	appToken := global.Config.Wxpusher.AppToken
	fmt.Println(appToken)
	data := struct {
		AppToken    string            `json:"appToken"`
		Content     string            `json:"content"`
		ContentType wxPushContentType `json:"contentType"`
		Summary     string            `json:"summary,omitempty"`
		TopicIds    []int             `json:"topicIds,omitempty"`
		Uids        []string          `json:"uids,omitempty"`
		Url         string            `json:"url,omitempty"`
	}{
		AppToken:    appToken,
		Content:     html,
		Summary:     title + ", " + desc,
		ContentType: p.contentType,
		TopicIds:    p.TopicIds,
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
