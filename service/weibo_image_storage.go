package service

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"hash/crc32"
	"strconv"
	"strings"
)

// 使用微博 实现图片云存储 微博的CDN实现加速
type WeiBoStorage struct {
	Sub string // cookie after login
}

func NewWeiBoStorage(SUB string) *WeiBoStorage {
	return &WeiBoStorage{Sub: SUB}
}

const (
	WeiBoUrl = "https://picupload.service.weibo.com/interface/pic_upload.php?mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog"
)

type UploadImageResponse struct {
	Code string                  `json:"code"`
	Data uploadImageResponseData `json:"data"`
}

type uploadImageResponseData struct {
	Data  string          `json:"data"`
	Count int             `json:"count"`
	Pics  map[string]pics `json:"pics"`
}

type pics struct {
	Width  int    `json:"width"`
	Size   int    `json:"size"`
	Ret    int    `json:"ret"`
	Height int    `json:"height"`
	Name   string `json:"name"`
	Pid    string `json:"pid"`
}

// 上传图片并返回可访问的地址
func (w *WeiBoStorage) UploadImage(imageData []byte) string {
	bs64data := base64.StdEncoding.EncodeToString(imageData)
	_, body, _ := gorequest.New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Post(WeiBoUrl).
		Type(gorequest.TypeMultipart).
		Set("Cookie", "SUB="+w.Sub).
		SendFile([]byte(bs64data), "b64_data", "b64_data").
		//SetDebug(conf.IsDev()).
		End()
	var t UploadImageResponse
	_ = json.Unmarshal(w.processResponseBody(body), &t)
	if t.Code == "A00006" {
		return w.pid2url(t.Data.Pics["pic_1"].Pid, "large")
	} else if t.Code == "A200001" {
		panic("WeiBo需要更新登录cookie")
	} else {
		panic("未知意义code：" + t.Code)
	}
}

func (w *WeiBoStorage) processResponseBody(html string) []byte {
	index := strings.Index(html, "{")
	jsonString := html[index:]
	return []byte(jsonString)
}

func (w *WeiBoStorage) pid2url(pid string, imageType string) string {
	var url string
	var zone uint32
	if pid[9] == 'w' {
		zone = (crc32.ChecksumIEEE([]byte(pid)) & 3) + 1
		url = fmt.Sprintf("http://ww%v.sinaimg.cn/%s/%s", zone, imageType, pid)
	} else {
		s := pid[len(pid)-2:]
		v, _ := strconv.ParseUint(s, 16, 32)
		zone = uint32((v & 0xf) + 1)
		url = fmt.Sprintf("http://ss%v.sinaimg.cn/%s/%s", zone, imageType, pid)
	}
	return url + ".jpg"
}
