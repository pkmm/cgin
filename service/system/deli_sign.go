package system

import (
	"cgin/global"
	"cgin/model/system"
	"cgin/util"
	"encoding/json"
	"errors"
	"github.com/parnurzeal/gorequest"
	"gorm.io/gorm"
	"net/http"
)

// TODO 重构 封装一个客户端

type DeliClient struct {
	userId string
	orgId  string
	token  string
}

func NewDeliClient() *DeliClient {
	return new(DeliClient)
}

func (d *DeliClient) GetUserIdAndToken(mobile, password string) (err error) {
	type loginData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	ver := loginData{mobile, util.MD5(password)}
	var result LoginRes

	// step 1, 拿到token
	_, _, errs := gorequest.New().Post("https://v2-app.delicloud.com/api/v2.0/auth/loginMobile").
		Set("client_id", "eplus_app").
		Set("X-Service-Id", "userauth").
		Set("User-Agent", "SmartOffice/2.4.4 (iPhone; iOS 15.0.1; Scale/3.00)").
		Send(ver).
		EndStruct(&result)

	if len(errs) > 0 {
		err = errs[0]
		return
	} else if result.Code != 0 {
		return errors.New(result.Msg)
	}

	if result.Code == 0 {
		d.userId = result.Data.UserId
		d.token = result.Data.Token
	}
	return err
}

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
		Set("client_type", "eplus_app").
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

func (d *DeliAutoSignService) SetAutoSign(username string, autoSign bool) (err error) {
	c := 1
	if autoSign {
		c = 0
	}
	err = global.DB.Model(system.DeliUser{}).Where("username = ?", username).Update("cancel", c).Error
	return err
}

func (d *DeliAutoSignService) Login(mobile, password string) (err error, token string) {
	type loginData struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	ver := loginData{mobile, util.MD5(password)}
	var result LoginRes

	// step 1, 拿到token
	_, _, errs := gorequest.New().Post("https://v2-app.delicloud.com/api/v2.0/auth/loginMobile").
		Set("client_id", "eplus_app").
		Set("X-Service-Id", "userauth").
		Set("User-Agent", "SmartOffice/2.4.4 (iPhone; iOS 15.0.1; Scale/3.00)").
		Send(ver).
		EndStruct(&result)

	if len(errs) > 0 {
		err = errs[0]
		return
	} else if result.Code != 0 {
		return errors.New(result.Msg), ""
	}

	// step 2, 拿到org_id
	type Step2RespData struct {
		Id string `json:"id"`
	}
	type Step2Resp struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data []Step2RespData `json:"data"`
	}
	var step2data Step2Resp
	rp, _, errs := gorequest.New().
		Get("https://v2-app.delicloud.com/api/v2.3/org/findOrgDetailByUserId?is_only_usable=true&user_id="+result.Data.UserId).
		Set("Authorization", result.Data.Token).
		Set("user_id", result.Data.UserId).
		Set("X-Service-Id", "organization").
		Set("client_id", "eplus_app").
		EndStruct(&step2data)
	if errs != nil {
		err = errs[0]
		return
	}
	if step2data.Code != 0 {
		err = errors.New(step2data.Msg)
		return
	}

	// step 3 拿到member_id
	type Step3RespData struct {
		OriginMemberId string `json:"origin_member_id"`
	}
	type Step3Resp struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data Step3RespData `json:"data"`
	}
	var step3data Step3Resp
	_, _, errs = gorequest.New().
		Get("https://v2-app.delicloud.com/api/v2.0/orgUser/findOrgUserDetailByOrgIdAndUserId?org_id="+step2data.Data[0].Id+"&user_id="+result.Data.UserId).
		Set("user_id", result.Data.UserId).
		Set("org_id", step2data.Data[0].Id).
		Set("X-Service-Id", "organization").
		Set("client_id", "eplus_app").
		Set("Authorization", result.Data.Token).
		Set("User-Agent", "SmartOffice/2.4.4 (iPhone; iOS 15.0.2; Scale/3.00)").
		EndStruct(&step3data)

	if errs != nil {
		err = errs[0]
		return
	}
	if step3data.Code != 0 {
		err = errors.New(step3data.Msg)
		return
	}

	// step 4，拿到cookie
	rp, _, errs = gorequest.New().
		Get("https://v2-kq.delicloud.com/attend/index/home").
		Set("client_type", "eplus_app").
		Set("token", result.Data.Token).
		Set("user_id", result.Data.UserId).
		Set("org_id", step2data.Data[0].Id).
		Set("v1_member_id", step3data.Data.OriginMemberId).
		Set("User-Agent", "SmartOffice/2.4.4 (iPhone; iOS 15.0.2; Scale/3.00)").
		End()
	if errs != nil {
		err = errs[0]
		return
	}

	for _, ck := range (*rp).Cookies() {
		if ck.Name == "deliUser" {
			token = ck.Value
			break
		}
	}

	return err, token
}
