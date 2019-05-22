package zcmu

import (
	"bytes"
	"cgin/util"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"image/gif"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 登录正方教务系统
// 密码错误五次后 当天是不能登录的
// 目前只支持查询成绩，其他的功能后期在做

const (
	host          = "zfxk.zcmu.edu.cn"
	baseUrl       = "https://" + host + "/"
	verifyCodeUrl = "CheckCode.aspx" // 验证码
	defaultUrl    = "default2.aspx"
	userAgent     = "Mozilla/5.0 (Windows NT 6.3; WOW64)" +
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36"

	POST      = "POST"
	GET       = "GET"
	viewState = "__VIEWSTATE"

	loginErrorMsgWrongPassword          = "密码错误！！"
	loginErrorMsgWrongVerifyCode        = "验证码不正确！！"
	loginErrorMsgCanNotLoginIn          = "密码错误，您密码输入错误已达5次，账号已锁定无法登录，次日自动解锁！如忘记密码，请与学院教学秘书联系!"
	loginErrorMsgWrongPasswordSometimes = "密码错误，您还有\\d次尝试机会！如忘记密码，请与学院教学秘书联系!"
	loginErrorMsgNotValidUser           = "用户名不存在或未按照要求参加教学活动！！"
)

var (
	CanNotGetViewStateException          = errors.New("无法解析出view state值")
	ClickQueryScoreButtonFailedException = errors.New("点击'查询成绩'按钮失败")
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
	client   *http.Client
	num      string
	pwd      string
	name     string
	loginUrl string // 提交数据登录页 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/default2.aspx ps: ()中的值可能存在
	baseURL  string // loginurl前半段 : => http://zfxk.zjtcm.net/(dbn5dgq4jveyap4525jo5j45)/
	errorMsg string // 登陆过程的错误信息

	scores []*Score //成绩结果

	viewState          string
	currentPageOfBytes []byte

	isPasswordWrong bool
}

func (c *Crawl) GetStudentName() string {
	return c.name
}

func (c *Crawl) GetErrorMsg() string {
	return c.errorMsg
}

func (c *Crawl) IsPassWordWrong() bool {
	return c.isPasswordWrong
}

// 是不是可以继续同步
func (c *Crawl) CanContinue() bool {
	return !c.isPasswordWrong && c.errorMsg != loginErrorMsgNotValidUser
}

// 检测账号的状态
func (c *Crawl) CheckAccount() (errorMsg string) {
	var utf8Html []byte
	var err error
	// 打开登录页
	if err = c.touchIndexPage(); err != nil {
		return err.Error()
	}
	// 提交表单登录
	if err = c.submitLoginForm(); err != nil {
		return err.Error()
	}
	// 检查登录后页面的状态值 确定是否登录成功
	if utf8Html, err = c.gbkToUtf8(c.currentPageOfBytes); err != nil {
		return err.Error()
	}
	regs := []*regexp.Regexp{
		regexp.MustCompile(loginErrorMsgWrongPassword),
		regexp.MustCompile(loginErrorMsgWrongVerifyCode),
		regexp.MustCompile(loginErrorMsgNotValidUser),
		regexp.MustCompile(loginErrorMsgCanNotLoginIn),
		regexp.MustCompile(loginErrorMsgWrongPasswordSometimes),
	}
	for _, reg := range regs {
		if value := reg.FindString(string(utf8Html)); value != "" {
			c.isPasswordWrong = strings.Index(value, "密码") != -1
			return value
		}
	}
	return ""
}

// 获取成绩的接口
func (c *Crawl) GetScores() ([]*Score, error) {

	if c.errorMsg = c.CheckAccount(); c.errorMsg != "" {
		return nil, errors.New(c.errorMsg)
	}

	if err := c.pressQueryScoreButton(); err != nil {
		return nil, err
	}

	if err := c.filterScores(); err != nil {
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
		return nil, errors.New("初始化cookie jar失败")
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
	return "", CanNotGetViewStateException
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
	ret := pattern.FindSubmatch(c.currentPageOfBytes)
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

// 第一步：访问系统首页 获取登陆的viewState
func (c *Crawl) touchIndexPage() error {
	var (
		resp *http.Response
		err  error
		html []byte
	)
	if resp, err = c.client.Get(baseUrl); err != nil {
		return err
	}
	defer resp.Body.Close()
	if html, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	c.currentPageOfBytes = html
	if c.viewState, err = c.getViewState(html); err != nil {
		return err
	}
	c.loginUrl = resp.Request.URL.String()

	// 服务端有可能开启 地址加入了(xxxxxxxxxx)这样的情况
	if c.loginUrl == baseUrl {
		c.baseURL = baseUrl
	} else {
		c.baseURL = c.loginUrl[:len(c.loginUrl)-len(defaultUrl)]
	}
	return nil
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

// 第二步：提交表单进行登录
func (c *Crawl) submitLoginForm() error {
	var (
		err  error
		code string
		rep  *http.Response
	)

	if code, err = c.verifyCode2String(); err != nil {
		return err
	}

	formData := url.Values{
		viewState:          {c.viewState},
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
	if c.currentPageOfBytes, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}
	return nil
}

// 第三步：点击`查询成绩`按钮 此时的页面view state需要保存，下一步看成绩需要使用
func (c *Crawl) pressQueryScoreButton() error {
	var (
		err  error
		req  *http.Request
		rep  *http.Response
		html []byte
	)

	// 获取 查询成绩 的按钮地址
	// (?=)正则表达式 顺序环视
	var re = regexp.MustCompile(`(?s)xscj_gc\.aspx\?xh=(.*?)\&xm=(.*?)\&gnmkdm=(.*?)"`)
	matches := re.FindAllStringSubmatch(string(c.currentPageOfBytes), -1)
	if matches == nil {
		return ClickQueryScoreButtonFailedException
	}

	c.name = matches[0][2] // "re", num, name, gnmkdm
	myUrl := fmt.Sprintf("%sxscj_gc.aspx?xh=%s&xm=%s&gnmkdm=N121605", c.baseURL, c.num, url.QueryEscape(c.name))
	if req, err = http.NewRequest(GET,
		myUrl,
		nil,
	); err != nil {
		return err
	}

	req.Header.Set("Referer", c.baseURL+"xs_main.aspx?xh="+c.num)

	if rep, err = c.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if html, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	c.currentPageOfBytes = html
	if c.viewState, err = c.getViewState(html); err != nil {
		return err
	}

	return nil
}

// 第四步：进入了成绩页面，选择filter进行显示成绩，这里暂时显示全部(//todo 使用filter)
func (c *Crawl) filterScores() error {
	var (
		err   error
		ddlXN = "" // 学年
		ddlXQ = "" // 学期
		req   *http.Request
		rep   *http.Response
	)

	formData := make(url.Values)
	formData.Set(viewState, c.viewState)
	formData.Set("ddlXN", ddlXN)
	formData.Set("ddlXQ", ddlXQ)
	formData.Set("Button2", "%D4%DA%D0%A3%D1%A7%CF%B0%B3%C9%BC%A8%B2%E9%D1%AF")

	if req, err = http.NewRequest(
		POST,
		c.baseURL+"xscj_gc.aspx?xh="+c.num+"&xm="+url.QueryEscape(c.name)+"&gnmkdm=N121605",
		strings.NewReader(formData.Encode()),
	); err != nil {
		fmt.Println(err, "sss")
		return err
	}

	req.Header.Set("Referer", c.baseURL+"/xs_main.aspx?xh="+c.num)
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	if rep, err = c.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if c.currentPageOfBytes, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	return nil
}
