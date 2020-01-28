package service

import (
	"bytes"
	"cgin/conf"
	"cgin/errno"
	"cgin/util"
	"encoding/json"
	"fmt"
	"golang.org/x/image/webp"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type dailyServ struct {
	baseService
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

type ImageTags struct {
	Name string `json:"name"`
}

type ImageUserProfilerImageUrl struct {
	Medium string `json:"medium"`
}

type ImageUser struct {
	Id               string                    `json:"id"`
	Name             string                    `json:"name"`
	Account          string                    `json:"account"`
	ProfileImageUrls ImageUserProfilerImageUrl `json:"profile_image_urls"`
}

type ImageItemMetaSinglePage struct {
	OriginalImageUrl string `json:"original_image_url,omitempty"`
	LargeImageUrl    string `json:"large_image_url,omitempty"`
}

type ImageURL struct {
	Original string `json:"original"`
	Large    string `json:"large"`
}

type ImageMetaPage struct {
	ImageUrls ImageURL `json:"image_urls"`
}

type ImageItem struct {
	Id             string                  `json:"id"`
	Title          string                  `json:"title"`
	Type           string                  `json:"type"`
	Caption        string                  `json:"caption"`
	User           ImageUser               `json:"user"`
	Tags           []ImageTags             `json:"tags"`
	Tools          []string                `json:"tools"`
	CreateDate     string                  `json:"create_date"`
	PageCount      int                     `json:"page_count"`
	Width          int                     `json:"width"`
	Height         int                     `json:"height"`
	SanityLevel    int                     `json:"sanity_level"`
	MetaSinglePage ImageItemMetaSinglePage `json:"meta_single_page"`
	MetaPages      []ImageMetaPage         `json:"meta_pages"`
}
type ImageAPIResponse struct {
	Message string      `json:"message"`
	Data    []ImageItem `json:"data,omitempty"`
}

func (d *dailyServ) getListByKeyword(page int, keyword string) *ImageAPIResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d&keyword=%s", conf.AppConfig.String("daily.image"), page, keyword), nil)
	d.CheckError(err)
	resp, err := client.Do(req)
	d.CheckError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	d.CheckError(err)
	result := &ImageAPIResponse{}
	err = json.Unmarshal(body, result)
	d.CheckError(err)
	return result
}

// 每日一图 使用
// https://sotama.cool/picture 接口
func (d *dailyServ) getImageFromAPI() (string ,error) {
	fileName := "static/images/" + util.Date() + ".jpg"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", conf.AppConfig.String("daily.image.api"), nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer out.Close()
	webpImage, err := webp.Decode(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	err = jpeg.Encode(out, webpImage, &jpeg.Options{Quality: 60})
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (d *dailyServ) GetImage() string {
	fileName := "static/images/" + util.Date() + ".jpg"
	// TODO: 逻辑是不是需要更新
	// 今天的图片已经存在那么就直接返回了
	if util.PathExists(fileName) {
		return fileName
	}

	getImageConfig := conf.AppConfig.DefaultString("daily.image.type", "not-api")
	if getImageConfig == "api" {
		if todayImageFilePath, err := d.getImageFromAPI(); err == nil {
			return todayImageFilePath
		}
	}

	page := rand.Intn(4) // TODO 页数的设置
	keyword := "saber"   // TODO 设置关键字 支持搜索
	list := d.getListByKeyword(page, keyword)
	itemIndex := rand.Intn(len(list.Data))
	item := list.Data[itemIndex]
	var imageUrl string
	var originalUrl string
	if item.PageCount > 1 {
		originalUrl = item.MetaPages[0].ImageUrls.Original
		imageUrl = strings.Replace(item.MetaPages[0].ImageUrls.Large, "_webp", "", -1)
	} else {
		originalUrl = item.MetaSinglePage.OriginalImageUrl
		imageUrl = strings.Replace(originalUrl, "img-original", "c/540x540_70/img-master", -1)
		imageUrl = strings.Replace(imageUrl, ".jpg", "_master1200.jpg", -1)
		imageUrl = strings.Replace(imageUrl, ".png", "_master1200.jpg", -1)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://img.pixivic.com:23334/get/%s", imageUrl), nil)
	d.CheckError(err)
	resp, err := client.Do(req)
	d.CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	d.CheckError(err)
	defer resp.Body.Close()
	out, err := os.Create(fileName)
	d.CheckError(err)
	defer out.Close()
	//webpImage, err := webp.Decode(bytes.NewReader(body))
	//d.CheckError(err)
	//err = jpeg.Encode(out, webpImage, &jpeg.Options{Quality: 60})
	//d.CheckError(err)
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
	_ = json.Unmarshal(body, &ret)
	// TODO:
	return ret
}
