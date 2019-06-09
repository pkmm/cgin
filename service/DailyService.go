package service

import (
	"bytes"
	"cgin/conf"
	"cgin/errno"
	"cgin/util"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type dailyServ struct {
	BaseService
}

var DailyService = &dailyServ{}

type DailySentenceResp struct {
	Id        int       `json:"id"`
	Hitokoto  string    `json:"hitokoto"`
	Type      string    `json:"type"`
	From      string    `json:"from"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *dailyServ) GetImage() string {
	fileName := "static/images/" + util.Date() + ".webp"
	// TODO: 逻辑是不是需要更新
	// 今天的图片已经存在那么就直接返回了
	if util.PathExists(fileName) {
		return fileName
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", conf.AppConfig.String("daily.image"), nil)
	d.CheckError(err)
	resp, err := client.Do(req)
	d.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	d.CheckError(err)
	defer resp.Body.Close()
	out, err := os.Create(fileName)
	d.CheckError(err)
	defer out.Close()
	_, err = io.Copy(out, bytes.NewReader(body))
	d.CheckError(err)
	// TODO: 保存在云，优化是不是需要保存下来
	return fileName
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
