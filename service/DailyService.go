package service

import (
	"cgin/conf"
	"cgin/errno"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type dailyServ struct {

}

var DailyService = &dailyServ{}

type DailySentenceResp struct{
	Id int `json:"id"`
	Hitokoto string `json:"hitokoto"`
	Type string `json:"type"`
	From string `json:"from"`
	Creator string `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *dailyServ) GetImage() (imageUrl string)  {
	client := &http.Client{}
	req, err := http.NewRequest("GET", conf.AppConfig.String("daily.image"), nil)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	defer resp.Body.Close()

	// TODO:
	return string(body)
}

func (d *dailyServ) GetSentence() (sentence *DailySentenceResp) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", conf.AppConfig.String("daily.sentence"), nil)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errno.NormalException.AppendErrorMsg(err.Error()))
	}
	defer resp.Body.Close()
	ret := &DailySentenceResp{}
	json.Unmarshal(body, &ret)
	// TODO:
	return ret
}