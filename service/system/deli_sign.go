package system

import (
	"cgin/model/system"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type DeliAutoSign struct {
}

var DeliAutoSignApp = new(DeliAutoSign)

type CheckResultUrl struct {
	Url string `json:"url"`
}
type CheckResult struct {
	Data   CheckResultUrl `json:"data"`
	Errno  int            `json:"errno"`
	Errmsg string         `json:"errmsg"`
}

func (d *DeliAutoSign) SignOne(user *system.DeliUser) (err error, html string) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 60 * time.Second}
	data := url.Values{
		"device_id":   {"210457"},
		"device_type": {"0"},
		"lat":         {"30.515479372560893"},
		"lng":         {"114.42014366321327"},
	}
	req, err := http.NewRequest(http.MethodPost,
		"https://v2-kq.delicloud.com/attend/check/check", strings.NewReader(data.Encode()))
	if err != nil {
		return err, ""
	}
	req.AddCookie(&http.Cookie{Name: "deliUser", Value: user.Token})
	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()
	var ret CheckResult
	json.Unmarshal(body, &ret)

	if resp, err = client.Get(ret.Data.Url); err != nil {
		return err, ""
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, ""
	}
	defer resp.Body.Close()
	return nil, string(body)
}
