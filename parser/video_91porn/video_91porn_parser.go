package video_91porn

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
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
		item.ViewKey = pElement.Find("a").First().AttrOr("href", "")
		videos = append(videos, item)
	})

	return &videos
}
