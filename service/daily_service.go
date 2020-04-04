package service

import (
	"bytes"
	"cgin/conf"
	"cgin/errno"
	"cgin/util"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"image"
	"image/jpeg"
	"os"
)

type DailySentenceResp struct {
	Id         int    `json:"id"`
	Hitokoto   string `json:"hitokoto"`
	Type       string `json:"type"`
	From       string `json:"from"`
	FromWho    string `json:"from_who,omitempty"`
	Creator    string `json:"creator"`
	CreatorUid int64  `json:"creator_uid"`
	Reviewer   int64  `json:"reviewer,omitempty"`
	UUID       string `json:"uuid"`
	CreatedAt  int64  `json:"created_at,string"`
}

type dailyServe struct {
	baseService
}

var DailyService = &dailyServe{}

// 返回图片的本地地址
func (d *dailyServe) getImageFromAPI(fileName string) error {
	var ret = struct {
		Error  float64 `json:"error,string"`
		Result float64 `json:"result,string"`
		Img    string  `json:"img"`
	}{}
	client := gorequest.New()
	client.Get("http://img.xjh.me/random_img.php?return=json").EndStruct(&ret)
	fmt.Printf("%#v", ret.Img)
	_, bodyBytes, errs := client.Get("http:" + ret.Img).EndBytes()
	if 0 != len(errs) {
		return errs[0]
	}

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	myImage, _, _ := image.Decode(bytes.NewReader(bodyBytes))
	err = jpeg.Encode(out, myImage, &jpeg.Options{Quality: 60})
	if err != nil {
		return err
	}
	return nil
}

func (d *dailyServe) GetImage() string {
	fileName := "static/images/" + util.Date() + ".jpg"
	if util.PathExists(fileName) {
		return fileName
	}
	if err := d.getImageFromAPI(fileName); err != nil {
		panic(errno.NormalException.ReplaceErrorMsgWith(err.Error()))
	}
	return fileName
}

func (d *dailyServe) GetSentence() (sentence *DailySentenceResp) {
	ret := &DailySentenceResp{}
	gorequest.New().Get(conf.AppConfig.String("daily.sentence")).EndStruct(ret)
	// TODO: 做自己的精彩的句子库
	return ret
}
