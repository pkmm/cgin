package system

import (
	"cgin/global"
	"cgin/model/system"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"net/http"
)

type DeliAutoSignService struct {
}

var DeliAutoSignApp = new(DeliAutoSignService)

type CheckResultUrl struct {
	Url string `json:"url"`
}
type CheckResult struct {
	Data   CheckResultUrl `json:"data"`
	Errno  int            `json:"errno"`
	Errmsg string         `json:"errmsg"`
}

func (d *DeliAutoSignService) SignOne(user *system.DeliUser) (err error, html string) {
	data := struct {
		DeviceId   int     `json:"device_id"`
		DeviceType int     `json:"device_type"`
		Lat        float64 `json:"lat"`
		Lng        float64 `json:"lng"`
	}{
		210457,
		0,
		30.515479372560893,
		114.42014366321327,
	}
	err = nil
	gorequest.New().Post("https://v2-kq.delicloud.com/attend/check/check").
		AddCookie(&http.Cookie{Name: "deliUser", Value: user.Token}).
		Send(data).
		End(func(response gorequest.Response, body string, errs []error) {
			var ret CheckResult
			err = json.Unmarshal([]byte(body), &ret)
			if err != nil {
				return
			}
			if ret.Data.Url != "" {
				_, html, errs = gorequest.New().Get(ret.Data.Url).End()
				if len(errs) > 0 {
					err = errs[0]
					html = "<body> ERROR </body>"
					return
				}
			}
		})
	return err, html
}

func (d *DeliAutoSignService) GetAllUsers() (err error, users []system.DeliUser) {
	err = global.DB.Where("cancel = 0").Find(&users).Error
	return
}

func (d *DeliAutoSignService) GetUserByName(name string) (user *system.DeliUser) {
	if err := global.DB.Where("username = ?", name).First(&user).Error; err != nil {
		return nil
	}
	return user
}

func (d *DeliAutoSignService) UpdateUserWxpushUID(username, uid string) (err error) {
	err = global.DB.Where("username = ?", username).Updates(&system.DeliUser{Uid: uid}).Error
	return err
}
