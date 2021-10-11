package request

type WxpushCbData struct {
	AppId       int    `json:"appId"`
	AppKey      string `json:"appKey"`
	AppName     string `json:"appName"`
	Source      string `json:"source"`
	Ussername   string `json:"ussername"`
	UserHeadImg string `json:"userHeadImg"`
	Time        int64  `json:"time"`
	Uid         string `json:"uid"`
	Extra       string `json:"extra"`
}

type WxpushCb struct {
	Action string `json:"action"`
	Data WxpushCbData `json:"data"`
}

