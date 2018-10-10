package utility

import (
	"bytes"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image/gif"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"pkmm/models"
	"regexp"
	"strings"
	"time"
)

// login zf system.
const (
	BASE_URL        = "http://zfxk.zjtcm.net/"
	LOGIN_URL       = "default2.aspx"
	VERIFY_CODE_URL = "CheckCode.aspx"

	POST      = "POST"
	GET       = "GET"
	VIEWSTATE = "__VIEWSTATE"
)

type Crawl struct {
	Client    *http.Client
	Num       string
	Pwd       string
	MainPage  string // 登陆成功后的mainPage页面 utf8编码
	ScorePage []byte // 成绩界面
}

func getViewState(html []byte) (string, error) {
	pattern := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}
	return "", errors.New("解析 viewstate 失败")
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (this *Crawl) retrieveScores() []*models.Score {
	// 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(this.ScorePage)
	var scores []*models.Score
	if len(ret) == 0 {
		return scores
	}
	table := ret[0]
	table, _ = GbkToUtf8(table)
	// <td>学年</td><td>学期</td><td>课程代码</td><td>课程名称</td><td>课程性质</td><td>课程归属</td><td>学分</td><td>绩点</td><td>成绩</td><td>辅修标记</td><td>补考成绩</td><td>重修成绩</td><td>学院名称</td><td>备注</td><td>重修标记</td><td>课程英文名称</td>
	pattern = regexp.MustCompile(`(?s)<td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>.*?</td><td>.*?</td><td>.*?</td>`)
	tds := pattern.FindAllSubmatch(table, -1)

	for index, row := range tds {
		if index == 0 {
			continue
		}
		if string(row[8]) == "&nbsp;" {
			row[8] = nil
		}
		if string(row[9]) == "&nbsp;" {
			row[9] = nil
		}
		score := &models.Score{
			Xn:   string(row[1]),
			Xq:   string(row[2]),
			Kcmc: string(row[3]),
			Type: string(row[4]),
			Xf:   string(row[5]),
			Jd:   string(row[6]),
			Cj:   string(row[7]),
			Bkcj: string(row[8]),
			Cxcj: string(row[9]),
		}
		scores = append(scores, score)
	}
	return scores
}

// 1. 打开登陆页
func (this *Crawl) openLoginPage() (string, error) {
	rep, err := this.Client.Get(BASE_URL)
	if err != nil {
		return "", errors.New("获取登陆页面失败")
	}
	defer rep.Body.Close()
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", errors.New("解析登陆页面失败")
	}
	viewState, err := getViewState(html)
	if err != nil {
		return "", errors.New("解析登录页的viewstate失败")
	}
	return viewState, nil
}

// 2. 获取验证码
func (this *Crawl) getCode() (string, error) {
	rep, err := this.Client.Get(BASE_URL + VERIFY_CODE_URL)
	if err != err {
		return "", errors.New("加载验证码失败")
	}
	defer rep.Body.Close()
	im, err := gif.Decode(rep.Body)
	if err != nil {
		return "", errors.New("解析验证码失败")
	}
	code, err := Predict(im, false)
	if err != nil {
		return "", err
	}
	return code, nil
}

// 3. 登陆后的主页
func (this *Crawl) GetMainPage() (string, error) {
	viewstate, err := this.openLoginPage()
	if err != nil {
		return "", err
	}
	code, err := this.getCode()
	if err != nil {
		return "", err
	}
	formData := url.Values{
		VIEWSTATE:          {viewstate},
		"txtUserName":      {this.Num},
		"Textbox1":         {""},
		"TextBox2":         {this.Pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	rep, err := this.Client.PostForm(BASE_URL+LOGIN_URL, formData)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}
	utf8Html, err := GbkToUtf8(html)
	if err != nil {
		return "", errors.New("转码登陆后的页面失败")
	}
	this.MainPage = string(utf8Html)
	return this.MainPage, nil
}

func (this *Crawl) regHelper(value string) (bool, string) {
	reg := regexp.MustCompile(value)
	if reg.FindString(this.MainPage) != "" {
		return false, value
	}
	return true, ""
}

func CheckAccount(num, pwd string) (bool, string) {
	var (
		err error
	)
	crawl := NewCrawl(num, pwd)
	if _, err = crawl.GetMainPage(); err != nil {
		return false, err.Error()
	}
	if ok, ret := crawl.regHelper("验证码不正确"); !ok {
		return ok, ret
	}

	if ok, ret := crawl.regHelper("密码错误"); !ok {
		return ok, ret
	}

	if ok, ret := crawl.regHelper("用户名不存在或未按照要求参加教学活动"); !ok {
		return ok, ret
	}

	if ok, ret := crawl.regHelper("用户名不存在或未按照要求参加教学活动"); !ok {
		return ok, ret
	}

	return true, ""
}

func NewCrawl(num, pwd string) *Crawl {
	timeout := time.Duration(3 * time.Second) // 超时3s
	crawl := &Crawl{}
	tmpJar, _ := cookiejar.New(nil)
	crawl.Client = &http.Client{
		Jar:     tmpJar,
		Timeout: timeout,
	}
	crawl.Num = num
	crawl.Pwd = pwd
	return crawl
}

func (this *Crawl) LoginScorePage() ([]*models.Score, error) {
	var (
		err          error
		scores       []*models.Score
		req          *http.Request
		rep          *http.Response
		newViewState string
	)

	this.GetMainPage()

	req, err = http.NewRequest(GET,
		"http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+this.Num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605",
		nil,
	)
	if err != nil {
		return scores, err
	}
	req.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+this.Num)
	if rep, err = this.Client.Do(req); err != nil {
		return scores, err
	}
	defer rep.Body.Close()

	if this.ScorePage, err = ioutil.ReadAll(rep.Body); err != nil {
		return scores, err
	}

	// 获取viewstate, 用于打开成绩页面
	if newViewState, err = getViewState(this.ScorePage); err != nil {
		return scores, err
	}

	var (
		ddlXN = ""
		ddlXQ = ""
	)

	formData := make(url.Values)
	formData.Set(VIEWSTATE, newViewState)
	formData.Set("ddlXN", ddlXN)
	formData.Set("ddlXQ", ddlXQ)
	formData.Set("Button2", "%D4%DA%D0%A3%D1%A7%CF%B0%B3%C9%BC%A8%B2%E9%D1%AF")

	req, err = http.NewRequest(
		POST,
		"http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+this.Num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return scores, err
	}
	req.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+this.Num)
	req.Header.Set("Host", "zfxk.zjtcm.net")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	rep, err = this.Client.Do(req)

	if err != nil {
		return scores, err
	}

	defer rep.Body.Close()

	this.ScorePage, err = ioutil.ReadAll(rep.Body)

	if err != nil {
		return scores, err
	}

	return this.retrieveScores(), nil
}
