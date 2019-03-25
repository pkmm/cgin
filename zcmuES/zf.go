package zcmuES

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
	"pkmm_gin/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 登录正方教务系统
// 密码错误五次后 当天是不能登录的
// 目前只支持查询成绩，其他的功能后期在做

const (
	baseUrl       = "http://zfxk.zjtcm.net/"
	verifyCodeUrl = "CheckCode.aspx" // 验证码
	defaultUrl    = "default2.aspx"
	host          = "zfxk.zjtcm.net"
	userAgent     = "Mozilla/5.0 (Windows NT 6.3; WOW64)" +
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36"

	POST      = "POST"
	GET       = "GET"
	viewState = "__VIEWSTATE"

	loginErrorMsgWrongPassword           = "密码错误！！"
	loginErrorMsgWrongVerifyCode         = "验证码不正确！！"
	loginErrorMsgCanNotLoginIn           = "密码错误，您密码输入错误已达5次，账号已锁定无法登录，次日自动解锁！如忘记密码，请与学院教学秘书联系!"
	loginErrorMsgNotValidUser            = "用户名不存在或未按照要求参加教学活动！！"
	loginErrorMsgDecodeViewStateError    = "解析viewstate失败"
	loginErrorMsgCanNotConnectZcmuSystem = "无法访问到教务系统"
	loginErrorMsgUnknown                 = "未知错误"
)

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

type Crawl struct {
	client                   *http.Client
	num                      string
	pwd                      string
	name                     string
	mainPage                 string // 登陆成功后的mainPage页面 utf8编码
	scorePageBeforeRealSeeIt []byte // 在进入成绩界面前 用于得到viewstate的页面
	scorePage                []byte // 点击了查看成绩的按钮之后 显示的页面， 这个页面就是包括了成绩的table的页面
	loginUrl                 string // 提交数据登录页 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/default2.aspx ps: ()中的值可能存在
	baseURL                  string // loginurl前半段 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/
	errorMsg                 string // 登陆过程的错误信息

	scores []*Score //成绩结果
}

func (c *Crawl) GetStudentName() string {
	return c.name
}

func (c *Crawl) GetErrorMsg() string {
	return c.errorMsg
}

func (c *Crawl) IsPassWordWrong() bool {
	return c.errorMsg == loginErrorMsgWrongPassword || c.errorMsg == loginErrorMsgCanNotLoginIn
}

// 是不是可以继续同步
func (c *Crawl) CanContinue() bool {
	return c.errorMsg != loginErrorMsgCanNotLoginIn && c.errorMsg != loginErrorMsgWrongPassword &&
		c.errorMsg != loginErrorMsgNotValidUser
}

// 检测账号的状态
func (c *Crawl) CheckAccount() (errorMsg string) {
	if err := c.prepareToLoginSystem(); err != nil {
		return err.Error()
	}

	html := c.mainPage
	regs := map[string]*regexp.Regexp{
		loginErrorMsgWrongPassword:   regexp.MustCompile(loginErrorMsgWrongPassword),
		loginErrorMsgWrongVerifyCode: regexp.MustCompile(loginErrorMsgWrongVerifyCode),
		loginErrorMsgNotValidUser:    regexp.MustCompile(loginErrorMsgNotValidUser),
		loginErrorMsgCanNotLoginIn:   regexp.MustCompile(loginErrorMsgCanNotLoginIn),
	}

	for key, reg := range regs {
		if reg.FindString(html) != "" {
			c.errorMsg = key
			return key
		}
	}

	return ""
}

// 获取成绩的接口
func (c *Crawl) GetScores() ([]*Score, error) {

	if msg := c.CheckAccount(); msg != "" {
		return nil, errors.New(msg)
	}

	if err := c.doLoginScorePage(); err != nil {
		return nil, err
	}

	return c.retrieveScores(), nil
}

func NewCrawl(num, pwd string) (*Crawl, error) {
	timeout := time.Duration(30 * time.Second) // 超时30s
	crawl := &Crawl{
		num: num,
		pwd: pwd,
	}
	tmpJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.New("初始化cookiejar失败")
	}
	crawl.client = &http.Client{
		Jar:     tmpJar,
		Timeout: timeout,
	}

	return crawl, nil
}

func (c *Crawl) getViewState(html []byte) (string, error) {
	pattern := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}

	return "", errors.New(loginErrorMsgDecodeViewStateError)
}

func (c *Crawl) gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (c *Crawl) retrieveScores() []*Score {
	// 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(c.scorePage)
	var scores []*Score
	if len(ret) == 0 {
		return scores
	}
	table := ret[0]
	table, _ = c.gbkToUtf8(table)
	// <td>学年</td><td>学期</td><td>课程代码</td><td>课程名称</td><td>课程性质</td><td>课程归属</td><td>学分</td><td>绩点</td><td>成绩</td><td>辅修标记</td><td>补考成绩</td><td>重修成绩</td><td>学院名称</td><td>备注</td><td>重修标记</td><td>课程英文名称</td>
	pattern = regexp.MustCompile(`(?s)<td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>.*?</td><td>.*?</td><td>.*?</td>`)
	tds := pattern.FindAllSubmatch(table, -1)

	for index, row := range tds {
		if index == 0 {
			continue
		}
		if string(row[8]) == "&nbsp;" {
			row[8] = []byte("")
		}
		if string(row[9]) == "&nbsp;" {
			row[9] = []byte("")
		}
		xq, _ := strconv.Atoi(string(row[2]))
		xf, _ := strconv.ParseFloat(string(row[5]), 64)
		jd, _ := strconv.ParseFloat(string(row[6]), 64)

		score := &Score{
			Xn:   string(row[1]),
			Xq:   uint8(xq),
			Kcmc: string(row[3]),
			Type: string(row[4]),
			Xf:   util.Decimal(xf),
			Jd:   util.Decimal(jd),
			Cj:   string(row[7]),
			Bkcj: string(row[8]),
			Cxcj: string(row[9]),
		}
		scores = append(scores, score)
	}
	return scores
}

// 未登录时, 访问系统首页 获取登陆的viewstate
func (c *Crawl) touchIndexPageForGetViewState() (string, error) {
	rep, err := c.client.Get(baseUrl)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}
	viewState, err := c.getViewState(html)
	if err != nil {
		return "", err
	}
	c.loginUrl = rep.Request.URL.String()

	// 服务端有可能开启 地址加入了(xxxxxxxxxx)这样的情况
	if c.loginUrl == baseUrl {
		c.baseURL = baseUrl
	} else {
		c.baseURL = c.loginUrl[:len(c.loginUrl)-len(defaultUrl)]
	}

	return viewState, nil
}

// 解析验证码
func (c *Crawl) verifyCode2String() (string, error) {
	rep, err := c.client.Get(c.baseURL + verifyCodeUrl)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()
	im, err := gif.Decode(rep.Body)
	if err != nil {
		return "", err
	}
	code, err := Predict(im)
	if err != nil {
		return "", err
	}

	return code, nil
}

// 尝试登陆系统, 获取登陆成功后的页面
func (c *Crawl) prepareToLoginSystem() error {
	var (
		err  error
		code string
		html []byte
		rep  *http.Response
		vs   string
	)
	if vs, err = c.touchIndexPageForGetViewState(); err != nil {
		return err
	}

	if code, err = c.verifyCode2String(); err != nil {
		return err
	}

	formData := url.Values{
		viewState:          {vs},
		"txtUserName":      {c.num},
		"Textbox1":         {""},
		"TextBox2":         {c.pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	request, _ := http.NewRequest(
		"POST",
		c.loginUrl,
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", c.loginUrl)
	request.Header.Set("Host", host)
	request.Header.Set("User-Agent", userAgent)

	if rep, err = c.client.Do(request); err != nil {
		return err
	}
	defer rep.Body.Close()

	if html, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}
	var utf8Html []byte
	if utf8Html, err = c.gbkToUtf8(html); err != nil {
		return err
	}

	c.mainPage = string(utf8Html)

	return nil
}

func (c *Crawl) touchScorePageForGetViewState() error {
	var (
		err error
		req *http.Request
		rep *http.Response
	)

	if err = c.prepareToLoginSystem(); err != nil {
		return err
	}

	// 获取 查询成绩 的按钮地址
	// (?=)正则表达式 顺序环视
	var re = regexp.MustCompile(`(?s)xscj_gc\.aspx\?xh=(.*?)\&xm=(.*?)\&gnmkdm=(.*?)"`)
	matches := re.FindAllStringSubmatch(c.mainPage, -1)
	if matches == nil {
		return errors.New("获取成绩按钮连接失败, num: [" + c.num + "], name: [" + c.name + "]")
	}

	c.name = matches[0][2] // "re", num, name, gnmkdm
	if req, err = http.NewRequest(GET,
		c.baseURL+matches[0][0][:len(matches[0][0])-1],
		nil,
	); err != nil {
		return err
	}

	req.Header.Set("Referer", c.baseURL+"xs_main.aspx?xh="+c.num)

	if rep, err = c.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if c.scorePageBeforeRealSeeIt, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	return nil
}

func (c *Crawl) doLoginScorePage() error {
	var (
		newViewState string
		err          error
		ddlXN        = ""
		ddlXQ        = ""
		req          *http.Request
		rep          *http.Response
	)

	if err = c.touchScorePageForGetViewState(); err != nil {
		return err
	}

	// 获取viewstate, 用于打开成绩页面
	if newViewState, err = c.getViewState(c.scorePageBeforeRealSeeIt); err != nil {
		return err
	}

	formData := make(url.Values)
	formData.Set(viewState, newViewState)
	formData.Set("ddlXN", ddlXN)
	formData.Set("ddlXQ", ddlXQ)
	formData.Set("Button2", "%D4%DA%D0%A3%D1%A7%CF%B0%B3%C9%BC%A8%B2%E9%D1%AF")

	if req, err = http.NewRequest(
		POST,
		c.baseURL+"xscj_gc.aspx?xh="+c.num+"&xm="+url.QueryEscape(c.name)+"&gnmkdm=N121605",
		strings.NewReader(formData.Encode()),
	); err != nil {
		return err
	}

	req.Header.Set("Referer", c.baseURL+"/xs_main.aspx?xh="+c.num)
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	if rep, err = c.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if c.scorePage, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	return nil
}
