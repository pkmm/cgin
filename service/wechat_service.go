package service

import (
	"bytes"
	"cgin/conf"
	"encoding/json"
	"fmt"
	"github.com/xlstudio/wxbizdatacrypt"
	"io/ioutil"
	"net/http"
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

// 模板消息 发送 接受的结构体
type TemplateMsg struct {
	AccessToken string           `json:"access_token"`
	ToUser      string           `json:"touser"`
	TemplateId  string           `json:"template_id"`
	Page        string           `json:"page"`
	//FormId      string           `json:"form_id"`
	Data        TemplateMsgData `json:"data"`
}

type TemplateMsgData struct {
	Key1 KeywordData `json:"thing1,omitempty"`
	Key2 KeywordData `json:"thing2,omitempty"`
}

type KeywordData struct {
	Value interface{} `json:"value"`
}

type SendTemplateResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	MsgID   string `json:"msgid"`
}

// 模板消息的部分结束

const (
	sendUserTemplateMsgUrl = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
	accessTokenUrl         = "https://api.weixin.qq.com/cgi-bin/token"
)

func getAccessToken() (token string, err error) {
	// todo token 缓存在redis中
	_url := fmt.Sprintf(accessTokenUrl+"?grant_type=client_credential&appid=%s&secret=%s",
		conf.AppConfig.String("miniprogram_app_id"),
		conf.AppConfig.String("miniprogram_secret"),
	)
	response, err := http.Get(_url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	type tmp struct {
		Token string `json:"access_token"`
	}
	t := new(tmp)
	json.Unmarshal(data, &t)
	return t.Token, nil
}

func SendUserTemplateMsg(msg *TemplateMsg) *SendTemplateResponse {
	token, err := getAccessToken()
	if err != nil {
		// todo
		return nil
	}
	msg.AccessToken = token
	data, err := json.Marshal(msg)
	if err != nil {
		// todo
		return nil
	}
	fmt.Printf("%#v", msg)
	resp, err := http.Post(fmt.Sprintf(sendUserTemplateMsgUrl, token), "application/json", bytes.NewBuffer(data))
	if err != nil {
		// todo
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	sendTemplateResponse := &SendTemplateResponse{}
	err = json.Unmarshal(body, sendTemplateResponse)
	if err != nil {
		// todo
		return nil
	}

	return sendTemplateResponse
}

func DecodeWeChatUserInfo(iv, code, encryptedData string) (*WechatUserInfo, error) {
	sess, err := code2Session(code)
	weChatUserInfo := &WechatUserInfo{}
	if err != nil {
		return weChatUserInfo, err
	}

	pc := wxbizdatacrypt.WxBizDataCrypt{
		AppID:      conf.AppConfig.String("miniprogram_app_id"),
		SessionKey: sess.SessionKey,
	}

	var result interface{}
	if result, err = pc.Decrypt(encryptedData, iv, true); err != nil {
		return weChatUserInfo, err
	}

	json.Unmarshal([]byte(result.(string)), &weChatUserInfo)

	return weChatUserInfo, nil
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
