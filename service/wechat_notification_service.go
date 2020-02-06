package service

import (
	"cgin/conf"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
)

type Notification interface {
	Notify(title, desc string) (*http.Response, error)
}

type weChatNotificationService struct {
	baseService
}

var WeChatNotificationService = &weChatNotificationService{}

type PushWorkerType int

const (
	ServerSister PushWorkerType = iota
	WxPusher
)

type wxPushContentType int

const (
	PlainText wxPushContentType = iota + 1
	Html
	Markdown
)

func (w *weChatNotificationService) Notify(
	title, desc string,
	workerType PushWorkerType,
	contentType wxPushContentType,
	uids []string,
) (*http.Response, error) {
	switch workerType {
	case ServerSister:
		return notifyWithServerSister(title, desc)
	case WxPusher:
		return notifyWithWxPusher(title, desc, contentType, uids)
	default:
		panic("cant not find push worker.")
	}
}

func notifyWithServerSister(title, desc string) (*http.Response, error) {
	key := conf.AppConfig.String("server_push_key")
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

type wxPushData struct {
	AppToken    string            `json:"appToken"`
	Content     string            `json:"content"`
	ContentType wxPushContentType `json:"contentType"`
	TopicIds    []int             `json:"topicIds,omitempty"`
	Uids        []string          `json:"uids,omitempty"`
	Url         string            `json:"url,omitempty"`
}

func notifyWithWxPusher(title, desc string, contentType wxPushContentType, uids []string) (*http.Response, error) {
	appToken := conf.AppConfig.String("app_token")
	data := &wxPushData{
		AppToken:    appToken,
		Content:     title + ", " + desc,
		ContentType: contentType,
		Uids:        uids,
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
