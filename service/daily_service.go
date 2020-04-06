package service

import (
	"bytes"
	"cgin/conf"
	"cgin/errno"
	"cgin/model"
	"cgin/util"
	"github.com/parnurzeal/gorequest"
	"image"
	"image/jpeg"
	"os"
	"sync"
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
	lock *sync.Mutex
}

var DailyService = &dailyServe{lock: &sync.Mutex{}}

// 返回图片的本地地址
// TODO: 使用微博作为图床，进行云存储, 本地可以不保存文件了
func (d *dailyServe) getImageFromAPI(fileName string, isSaveLocal bool, isSaveCloud bool) []error {
	var ret = struct {
		Error  float64 `json:"error,string"`
		Result float64 `json:"result,string"`
		Img    string  `json:"img"`
	}{}
	client := gorequest.New()
	client.Get("http://img.xjh.me/random_img.php?return=json").EndStruct(&ret)
	_, bodyBytes, errs := client.Get("http:" + ret.Img).EndBytes()
	if 0 != len(errs) {
		return errs
	}

	if isSaveLocal {
		errs = append(errs, saveImageLocal(bodyBytes, fileName))
	}

	if isSaveCloud {
		errs = append(errs, saveImageCloud(bodyBytes))
	}

	return nil
}

func saveImageCloud(data []byte) error {
	path := NewWeiBoStorage(conf.WeiBoCookie()).UploadImage(data)
	m := &model.ImageStorageInfo{Url: path}
	_, err := m.Create()
	return err
}

func saveImageLocal(data []byte, fileName string) error {
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	myImage, _, _ := image.Decode(bytes.NewReader(data))
	err = jpeg.Encode(out, myImage, &jpeg.Options{Quality: 60})
	if err != nil {
		return err
	}
	return nil
}

// 返回图片的资源地址
func (d *dailyServe) GetImage() string {
	defer func() {
		d.lock.Unlock()
	}()
	d.lock.Lock()
	imageModel := &model.ImageStorageInfo{}
	info, ok := imageModel.FindTodayImage()
	if ok && info != nil {
		return info.Url
	}
	fileName := "static/images/" + util.Date() + ".jpg"
	if util.PathExists(fileName) {
		return conf.Host() + "/" + fileName
	}
	if errs := d.getImageFromAPI(fileName, true, true); len(errs) != 0 {
		panic(errno.NormalException.ReplaceErrorByErrors(errs))
	}
	return conf.Host() + "/" + fileName
}

func (d *dailyServe) GetSentence() (sentence *DailySentenceResp) {
	ret := &DailySentenceResp{}
	gorequest.New().Get(conf.AppConfig.String("daily.sentence")).EndStruct(ret)
	// TODO: 做自己的精彩的句子库
	return ret
}
