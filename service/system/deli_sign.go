package system

import (
	"cgin/global"
	"cgin/model/system"
	"cgin/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gorm.io/gorm"
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

type LoginResultData struct {
	Expire string `json:"expire"`
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type LoginRes struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data LoginResultData `json:"data"`
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

func (d *DeliAutoSignService) UpdateUserToken(username, token string) (err error) {
	err = global.DB.Where("username = ?", username).Updates(&system.DeliUser{Token: token}).Error
	return err
}

func (d *DeliAutoSignService) UserExists(username string) (ok bool) {
	err := global.DB.Where("username = ?", username).First(&system.DeliUser{}).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (d *DeliAutoSignService) CreateUser(username, token string, useAutoSign bool) (err error) {
	cancel := 1 // cancel == 0 表示使用自动签到
	if useAutoSign {
		cancel = 0
	}
	err = global.DB.Create(&system.DeliUser{Username: username, Token: token, Cancel: cancel}).Error
	return err
}

func (d *DeliAutoSignService) Login(mobile, password string) (err error, result LoginRes) {
	type loginData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	ver := loginData{mobile, util.MD5(password)}

	_, _, errs := gorequest.New().Post("https://v2-app.delicloud.com/api/v2.0/auth/loginMobile").
		Set("client_id", "eplus_app").
		Set("X-Service-Id", "userauth").
		Set("User-Agent", "SmartOffice/2.4.4 (iPhone; iOS 15.0.1; Scale/3.00)").
		Send(ver).EndStruct(&result)
	fmt.Printf("data： %v\n", result)
	if len(errs) > 0 {
		err = errs[0]
		return
	} else if result.Code != 0 {
		return errors.New(result.Msg), result
	}
	return err, result
}
