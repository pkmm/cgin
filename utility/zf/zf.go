package zf

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
	"pkmm_gin/utility"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 登录正方教务系统
// 密码错误五次后 当天是不能登录的
// 目前只支持查询成绩，其他的功能后期在做

const (
	BASE_URL        = "http://zfxk.zjtcm.net/"
	VERIFY_CODE_URL = "CheckCode.aspx" // 验证码
	DEFAULT_URL     = "default2.aspx"
	HOST            = "zfxk.zjtcm.net"
	USER_ARENT      = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36"

	POST      = "POST"
	GET       = "GET"
	VIEWSTATE = "__VIEWSTATE"

	LOGIN_ERROR_MSG_WRONG_PASSWORD          = "密码错误"
	LOGIN_ERROR_MSG_WRONG_CHECK_CODE        = "验证码不正确！！"
	LOGIN_ERROR_MSG_NOT_VALID_USER          = "用户名不存在或未按照要求参加教学活动"
	LOGIN_ERROR_MSG_DECODE_VIEW_STATE_ERROR = "解析viewstate失败"

	LOGIN_ERROR_CODE_WRONG_PASSWORD          = 1001
	LOGIN_ERROR_CODE_WRONG_VERIFY_CODE       = 1002
	LOGIN_ERROR_CODE_NOT_VALID_USER          = 1003
	LOGIN_ERROR_CODE_DECODE_VIEW_STATE_ERROR = 1004
)

type Crawl struct {
	Client    *http.Client
	Num       string
	Pwd       string
	Name      string
	MainPage  string // 登陆成功后的mainPage页面 utf8编码
	ScorePage []byte // 成绩界面
	LoginUrl  string // 提交数据登录页 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/default2.aspx ps: ()中的值可能存在
	BaseURL   string // loginurl前半段 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/
}

func NewCrawl(num, pwd string) *Crawl {
	timeout := time.Duration(30 * time.Second) // 超时30s
	crawl := &Crawl{}
	tmpJar, err := cookiejar.New(nil)
	if err != nil {
		panic("初始化cookiejar失败")
	}
	crawl.Client = &http.Client{
		Jar:     tmpJar,
		Timeout: timeout,
	}
	crawl.Num = num
	crawl.Pwd = pwd
	return crawl
}

// 成绩结构
type Score struct {
	Xn   string  // 学年
	Xq   uint8   // 学期
	Kcmc string  // 课程名称
	Type string  // 课程性质
	Xf   float64 // 学分
	Jd   float64 // 绩点
	Cj   string  // 成绩
	Bkcj string  // 补考成绩
	Cxcj string  // 重修成绩
}

func getViewState(html []byte) (string, error) {
	pattern := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}
	return "", errors.New(LOGIN_ERROR_MSG_DECODE_VIEW_STATE_ERROR)
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (this *Crawl) retrieveScores() []*Score {
	// 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(this.ScorePage)
	var scores []*Score
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
		xq, _ := strconv.Atoi(string(row[2]))
		xf, _ := strconv.ParseFloat(string(row[5]), 64)
		jd, _ := strconv.ParseFloat(string(row[6]), 64)

		score := &Score{
			Xn:   string(row[1]),
			Xq:   uint8(xq),
			Kcmc: string(row[3]),
			Type: string(row[4]),
			Xf:   utility.Decimal(xf),
			Jd:   utility.Decimal(jd),
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
	this.LoginUrl = rep.Request.URL.String()

	// 服务端有可能开启 地址加入了(xxxxxxxxxx)这样的情况
	if this.LoginUrl == BASE_URL {
		this.BaseURL = BASE_URL
	} else {
		this.BaseURL = this.LoginUrl[:len(this.LoginUrl)-len(DEFAULT_URL)]
	}
	//fmt.Println(this.LoginUrl, this.BaseURL)
	return viewState, nil
}

// 2. 获取验证码
func (this *Crawl) getCode() (string, error) {
	rep, err := this.Client.Get(this.BaseURL + VERIFY_CODE_URL)
	if err != nil {
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
	request, _ := http.NewRequest(
		"POST",
		this.LoginUrl,
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", this.LoginUrl)
	request.Header.Set("Host", HOST)
	request.Header.Set("User-Agent", USER_ARENT)

	rep, err := this.Client.Do(request)
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

func (this *Crawl) checkLoginResult() (ok bool, value string) {
	if ok, ret := this.regHelper(LOGIN_ERROR_MSG_WRONG_CHECK_CODE); !ok {
		return ok, ret
	}

	if ok, ret := this.regHelper(LOGIN_ERROR_MSG_WRONG_PASSWORD); !ok {
		return ok, ret
	}

	if ok, ret := this.regHelper(LOGIN_ERROR_MSG_NOT_VALID_USER); !ok {
		return ok, ret
	}
	return true, ""
}

func (this *Crawl) regHelper(value string) (bool, string) {
	reg := regexp.MustCompile(value)
	if reg.FindString(this.MainPage) != "" {
		return false, value
	}
	return true, ""
}

func (this *Crawl) LoginScorePage() ([]*Score, error) {
	var (
		err          error
		scores       []*Score
		req          *http.Request
		rep          *http.Response
		newViewState string
	)

	this.GetMainPage()
	if ok, msg := this.checkLoginResult(); !ok {
		return scores, errors.New(msg)
	}
	// 获取 查询成绩 的按钮地址
	// (?=)正则表达式 顺序环视
	var re = regexp.MustCompile(`(?s)xscj_gc\.aspx\?xh=(.*?)\&xm=(.*?)\&gnmkdm=(.*?)"`)
	matches := re.FindAllStringSubmatch(this.MainPage, -1)
	if matches == nil {
		return scores, errors.New("获取成绩按钮连接失败" + this.Num + this.Name)
	}

	this.Name = matches[0][2] // "re", num, name, gnmkdm
	req, err = http.NewRequest(GET,
		this.BaseURL+matches[0][0][:len(matches[0][0])-1],
		nil,
	)
	if err != nil {
		return scores, err
	}
	req.Header.Set("Referer", this.BaseURL+"xs_main.aspx?xh="+this.Num)
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
		this.BaseURL+"xscj_gc.aspx?xh="+this.Num+"&xm="+url.QueryEscape(this.Name)+"&gnmkdm=N121605",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return scores, err
	}
	req.Header.Set("Referer", this.BaseURL+"/xs_main.aspx?xh="+this.Num)
	req.Header.Set("Host", HOST)
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
