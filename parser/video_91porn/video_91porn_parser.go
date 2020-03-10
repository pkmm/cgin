package video_91porn

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HOST      = "http://91porn.com/"
	Index     = "index.php"
	proxyHost = "http://localhost:1081"
)

type videoItem struct {
	Id            int64
	ViewKey       string
	Title         string
	ImgUrl        string
	Duration      string
	Info          string
	VideoResultId int64
}

type MyClient struct {
	c *http.Client
}

func (m *MyClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8") // 保证返回是中文
	//req.Header.Add("")
	return m.c.Do(req)
}

func setProxy(proxyHost string) (*url.URL, error) {
	u, err := url.Parse(proxyHost)
	return u, err
}

func NewHttpClient() *MyClient {
	proxy, err := setProxy(proxyHost)
	if err != nil {
		log.Fatal(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	return &MyClient{c: client}
}

// 解析主页视频
func ParseIndex() *[]videoItem {

	client := NewHttpClient()
	request, err := http.NewRequest("GET", HOST+Index, nil)
	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	videos := make([]videoItem, 0)
	doc.Find("div#tab-featured p").Each(func(i int, pElement *goquery.Selection) {
		item := videoItem{}
		item.Title = pElement.Find(".title").Text()
		item.ImgUrl = pElement.Find("img").First().AttrOr("src", "")
		item.Duration = pElement.Find(".duration").First().Text()
		href := pElement.Find("a").First().AttrOr("href", "")
		item.ViewKey = href[strings.LastIndex(href, "=")+1:]
		//allInfo := pElement.Text()
		//start := strings.Index(allInfo, "添加时间")
		//item.Info = allInfo[start:]

		videos = append(videos, item)
	})

	return &videos
}


// 解析其他的类别
func parseByCategory(html io.Reader) *[]videoItem {
	videoItemList := make([]videoItem, 0)
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("#fullside").Find(".listchannel").Each(func(i int, element *goquery.Selection) {
		item := videoItem{}
		contentUrl := element.Find("a").First().AttrOr("href", "")
		contentUrl = contentUrl[0: strings.Index(contentUrl, "&")]
		item.ViewKey = contentUrl[strings.Index(contentUrl, "=")+1:]
		tmpElement := element.Find("a").First().Find("img").First()
		item.ImgUrl = tmpElement.AttrOr("src", "")
		item.Title = tmpElement.AttrOr("title", "")
		allInfo := element.Text()
		startIndex := strings.Index(allInfo, "时长")
		item.Duration = allInfo[startIndex+6: startIndex+12]
		//start := strings.Index(allInfo, "添加时间")
		//info := allInfo[start:]
		//item.Info = strings.Replace(info, "还未被评分", "", -1)
		videoItemList = append(videoItemList, item)
	})
	return &videoItemList
}

//// 解析视屏的播放链接
//func ParseVideoPlayUrl(html io.Reader) videoItem {
//
//}