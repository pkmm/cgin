package util

import (
	"encoding/json"
	"fmt"
	"github.com/xlstudio/wxbizdatacrypt"
	"io/ioutil"
	"net/http"
	"pkmm_gin/conf"
)

// 使用code换取sessionKey
type code2SessionResp struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
	Errcode    string `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}

type WechatUserInfo struct {
	OpenId    string `json:"OpenId"`
	NickName  string `json:"nickName"`
	Language  string `json:"language"`
	Province  string `json:"province"`
	AvatarUrl string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
}

func DecodeWchatUserInfo(iv, code, encryptedData string) (*WechatUserInfo, error) {
	sess, err := code2Session(code)
	wechatUserInfo := &WechatUserInfo{}
	if err != nil {
		return wechatUserInfo, err
	}

	pc := wxbizdatacrypt.WxBizDataCrypt{
		AppID:      conf.AppConfig.String("miniprogram_app_id"),
		SessionKey: sess.SessionKey,
	}

	var result interface{}
	if result, err = pc.Decrypt(encryptedData, iv, true); err != nil {
		return wechatUserInfo, err
	}

	json.Unmarshal([]byte(result.(string)), &wechatUserInfo)

	return wechatUserInfo, nil
}

func code2Session(code string) (*code2SessionResp, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appId=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		conf.AppConfig.String("miniprogram_app_id"),
		conf.AppConfig.String("miniprogram_secret"),
		code,
	)

	sess := &code2SessionResp{}
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return sess, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sess, err
	}

	json.Unmarshal(body, &sess)

	return sess, nil
}
