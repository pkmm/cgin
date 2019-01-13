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
	loginErrorMsgWrongCheckCode          = "验证码不正确！！"
	loginErrorMsgCanNotLoginIn           = "密码错误，您密码输入错误已达5次，账号已锁定无法登录，次日自动解锁！如忘记密码，请与学院教学秘书联系!"
	loginErrorMsgNotValidUser            = "用户名不存在或未按照要求参加教学活动！！"
	loginErrorMsgDecodeViewStateError    = "解析viewstate失败"
	loginErrorMsgCanNotConnectZcmuSystem = "无法访问到教务系统"
	loginErrorMsgUnknown                 = "未知错误"

	loginErrorCodeWrongPassword        = 1001
	loginErrorCodeWrongVerifyCode      = 1002
	loginErrorCodeNotValidUser         = 1003
	loginErrorCodeDecodeViewStateError = 1004
	loginErrorCanNotOpenZcmuSystem     = 1005
	loginErrorCanNotLoginIn            = 1006
)

// 成绩结构
type score struct {
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
	errorCode                int    // 登陆过程的错误代码

	scores []*score //成绩结果
}

func (this *Crawl) GetErrorCode() int {
	return this.errorCode
}

func (this *Crawl) GetErrorMsg() string {
	return this.mapCode2Msg(this.errorCode)
}

func (this *Crawl) IsPassWordWrong() bool {
	return this.errorCode == loginErrorCodeWrongPassword
}

// 是不是可以继续同步
func (this *Crawl) CanContinue() bool {
	return this.errorCode != loginErrorCodeWrongPassword && this.errorCode != loginErrorCodeNotValidUser &&
		this.errorCode != loginErrorCanNotLoginIn
}

// 检测账号的状态
func (this *Crawl) CheckAccount() (errorMsg string) {
	if err := this.prepareToLoginSystem(); err != nil {
		return err.Error()
	}

	html := this.mainPage
	regs := map[string]*regexp.Regexp{
		loginErrorMsgWrongPassword:  regexp.MustCompile(loginErrorMsgWrongPassword),
		loginErrorMsgWrongCheckCode: regexp.MustCompile(loginErrorMsgWrongCheckCode),
		loginErrorMsgNotValidUser:   regexp.MustCompile(loginErrorMsgNotValidUser),
		loginErrorMsgCanNotLoginIn:  regexp.MustCompile(loginErrorMsgCanNotLoginIn),
	}

	for key, reg := range regs {
		if reg.FindString(html) != "" {
			this.errorCode = this.convertMsg2Code(key)
			return key
		}
	}

	return ""
}

func (this *Crawl) convertMsg2Code(msg string) int {
	switch msg {
	case loginErrorMsgWrongPassword:
		return loginErrorCodeWrongPassword
	case loginErrorMsgWrongCheckCode:
		return loginErrorCodeWrongVerifyCode
	case loginErrorMsgNotValidUser:
		return loginErrorCodeNotValidUser
	case loginErrorMsgDecodeViewStateError:
		return loginErrorCodeDecodeViewStateError
	case loginErrorMsgCanNotLoginIn:
		return loginErrorCanNotLoginIn
	default:
		return loginErrorCanNotOpenZcmuSystem
	}
}

// 获取成绩的接口
func (this *Crawl) GetScores() ([]*score, error) {

	if msg := this.CheckAccount(); msg != "" {
		return nil, errors.New(msg)
	}

	if err := this.doLoginScorePage(); err != nil {
		return nil, err
	}

	return this.retrieveScores(), nil
}

func (this *Crawl) mapCode2Msg(code int) string {
	switch code {
	case loginErrorCodeWrongPassword:
		return loginErrorMsgWrongPassword
	case loginErrorCodeWrongVerifyCode:
		return loginErrorMsgWrongCheckCode
	case loginErrorCodeNotValidUser:
		return loginErrorMsgNotValidUser
	case loginErrorCodeDecodeViewStateError:
		return loginErrorMsgDecodeViewStateError
	case loginErrorCanNotOpenZcmuSystem:
		return loginErrorMsgCanNotConnectZcmuSystem
	default:
		return loginErrorMsgUnknown
	}
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

func (this *Crawl) getViewState(html []byte) (string, error) {
	pattern := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}

	return "", errors.New(loginErrorMsgDecodeViewStateError)
}

func (this *Crawl) gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (this *Crawl) retrieveScores() []*score {
	// 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(this.scorePage)
	var scores []*score
	if len(ret) == 0 {
		return scores
	}
	table := ret[0]
	table, _ = this.gbkToUtf8(table)
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

		score := &score{
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
func (this *Crawl) touchIndexPageForGetViewState() (string, error) {
	rep, err := this.client.Get(baseUrl)
	if err != nil {
		return "", err
	}
	defer rep.Body.Close()
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}
	viewState, err := this.getViewState(html)
	if err != nil {
		return "", err
	}
	this.loginUrl = rep.Request.URL.String()

	// 服务端有可能开启 地址加入了(xxxxxxxxxx)这样的情况
	if this.loginUrl == baseUrl {
		this.baseURL = baseUrl
	} else {
		this.baseURL = this.loginUrl[:len(this.loginUrl)-len(defaultUrl)]
	}

	return viewState, nil
}

// 解析验证码
func (this *Crawl) verifyCode2String() (string, error) {
	rep, err := this.client.Get(this.baseURL + verifyCodeUrl)
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
func (this *Crawl) prepareToLoginSystem() error {
	var (
		err  error
		code string
		html []byte
		rep  *http.Response
		vs   string
	)
	if vs, err = this.touchIndexPageForGetViewState(); err != nil {
		return err
	}

	if code, err = this.verifyCode2String(); err != nil {
		return err
	}

	formData := url.Values{
		viewState:          {vs},
		"txtUserName":      {this.num},
		"Textbox1":         {""},
		"TextBox2":         {this.pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	request, _ := http.NewRequest(
		"POST",
		this.loginUrl,
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", this.loginUrl)
	request.Header.Set("Host", host)
	request.Header.Set("User-Agent", userAgent)

	if rep, err = this.client.Do(request); err != nil {
		return err
	}
	defer rep.Body.Close()

	if html, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}
	var utf8Html []byte
	if utf8Html, err = this.gbkToUtf8(html); err != nil {
		return err
	}

	this.mainPage = string(utf8Html)

	return nil
}

func (this *Crawl) touchScorePageForGetViewState() error {
	var (
		err error
		req *http.Request
		rep *http.Response
	)

	if err = this.prepareToLoginSystem(); err != nil {
		return err
	}

	// 获取 查询成绩 的按钮地址
	// (?=)正则表达式 顺序环视
	var re = regexp.MustCompile(`(?s)xscj_gc\.aspx\?xh=(.*?)\&xm=(.*?)\&gnmkdm=(.*?)"`)
	matches := re.FindAllStringSubmatch(this.mainPage, -1)
	if matches == nil {
		return errors.New("获取成绩按钮连接失败, num: [" + this.num + "], name: [" + this.name + "]")
	}

	this.name = matches[0][2] // "re", num, name, gnmkdm
	if req, err = http.NewRequest(GET,
		this.baseURL+matches[0][0][:len(matches[0][0])-1],
		nil,
	); err != nil {
		return err
	}

	req.Header.Set("Referer", this.baseURL+"xs_main.aspx?xh="+this.num)

	if rep, err = this.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if this.scorePageBeforeRealSeeIt, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	return nil
}

func (this *Crawl) doLoginScorePage() error {
	var (
		newViewState string
		err          error
		ddlXN        = ""
		ddlXQ        = ""
		req          *http.Request
		rep          *http.Response
	)

	if err = this.touchScorePageForGetViewState(); err != nil {
		return err
	}

	// 获取viewstate, 用于打开成绩页面
	if newViewState, err = this.getViewState(this.scorePageBeforeRealSeeIt); err != nil {
		return err
	}

	formData := make(url.Values)
	formData.Set(viewState, newViewState)
	formData.Set("ddlXN", ddlXN)
	formData.Set("ddlXQ", ddlXQ)
	formData.Set("Button2", "%D4%DA%D0%A3%D1%A7%CF%B0%B3%C9%BC%A8%B2%E9%D1%AF")

	if req, err = http.NewRequest(
		POST,
		this.baseURL+"xscj_gc.aspx?xh="+this.num+"&xm="+url.QueryEscape(this.name)+"&gnmkdm=N121605",
		strings.NewReader(formData.Encode()),
	); err != nil {
		return err
	}

	req.Header.Set("Referer", this.baseURL+"/xs_main.aspx?xh="+this.num)
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	if rep, err = this.client.Do(req); err != nil {
		return err
	}

	defer rep.Body.Close()

	if this.scorePage, err = ioutil.ReadAll(rep.Body); err != nil {
		return err
	}

	return nil
}
