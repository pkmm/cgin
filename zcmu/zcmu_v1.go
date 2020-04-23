package zcmu

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// use:  http://xpj.zcmu.edu.cn/jwglxt/xtgl/login_slogin.html?time=1587032918089

type ZSpider struct {
	number, pwd string
	client      *http.Client
	indexUrl    string
	errs        []error
}

const (
	rsaParams = "http://xpj.zcmu.edu.cn/jwglxt/xtgl/login_getPublicKey.html"
)

var (
	ErrCannotOpenPage            = errors.New("无法打开页面")
	ErrGetCSRFTokenFailed        = errors.New("获取csrfToken失败")
	ErrGenerateRsaPasswordFailed = errors.New("生成pwd的rsa失败")
	ErrGetRsaModuleFailed        = errors.New("获取rsa的module和exponent失败")
	ErrUsernameOrPasswordWrong   = errors.New("用户名或者密码错误")
	ErrLoginInOther              = errors.New("您的账号在其它地方登录，您已被迫下线。若非本人操作，请及时修改密码！")
	loginUrl                     = "http://xpj.zcmu.edu.cn/jwglxt/xtgl/login_slogin.html?time="
	kcUrl                        = "http://xpj.zcmu.edu.cn/jwglxt/cjcx/cjcx_cxDgXscj.html?doType=query&gnmkdm=N100801&su=%s"
)

func New(number string, pwd string) *ZSpider {
	// load cookie.
	jar, _ := cookiejar.New(nil)
	return &ZSpider{number: number, pwd: pwd, client: &http.Client{Jar: jar}}
}

func (z *ZSpider) getCSRFToken() (string, error) {
	z.indexUrl = fmt.Sprintf("http://xpj.zcmu.edu.cn/jwglxt/xtgl/login_slogin.html?time=%d", time.Now().UnixNano()/1e6)
	resp, err := z.client.Get(z.indexUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", ErrGetCSRFTokenFailed
	}
	elem := doc.Find("#csrftoken").AttrOr("value", "")
	return elem, nil
}

func (z *ZSpider) getRsaParams() (error, string, string) {
	n := time.Now().UnixNano() / 1e6
	resp, err := z.client.Get(fmt.Sprintf("%s?time=%d&_=%d", rsaParams, n, n-10))
	if err != nil {
		return err, "", ""
	}
	defer resp.Body.Close()
	bytesData, _ := ioutil.ReadAll(resp.Body)

	var ret = struct {
		Modulus  string `json:"modulus"`
		Exponent string `json:"exponent"`
	}{}
	_ = json.Unmarshal(bytesData, &ret)
	return nil, ret.Modulus, ret.Exponent
}

func (z *ZSpider) encryptPwd() (string, error) {
	err, module, exponent := z.getRsaParams()
	if err != nil {
		return "", err
	}
	var mod, ex []byte
	mod, _ = base64.StdEncoding.DecodeString(module)
	ex, _ = base64.StdEncoding.DecodeString(exponent)
	pub := &rsa.PublicKey{
		N: new(big.Int).SetBytes(mod),
		E: int(new(big.Int).SetBytes(ex).Int64()),
	}
	ciphered, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(z.pwd))
	if err != nil {
		return "", ErrGenerateRsaPasswordFailed
	}
	mm := base64.StdEncoding.EncodeToString(ciphered)
	return mm, nil
}

func (z *ZSpider) Login() error {
	token, err := z.getCSRFToken()
	if err != nil {
		return err
	}

	mm, err := z.encryptPwd()

	url.Values{}.Encode()
	// 基础知识要牢固 url query encode
	sdata := fmt.Sprintf("csrftoken=%s&mm=%s&yhm=%s&mm=%s", url.QueryEscape(token), url.QueryEscape(mm), url.QueryEscape(z.number), url.QueryEscape(mm))
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%d", loginUrl, time.Now().UnixNano()/1e6), strings.NewReader(sdata))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	req.Header.Set("Referer", "http://xpj.zcmu.edu.cn/jwglxt/xtgl/login_slogin.html")
	req.Header.Set("Host", "xpj.zcmu.edu.cn")
	req.Header.Set("Origin", "http://xpj.zcmu.edu.cn")
	req.Header.Set("User-Agent", userAgent)

	resp, err := z.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytesData, _ := ioutil.ReadAll(resp.Body)
	// TODO: 有时间进行优化 字符串查找
	ind := bytes.Index(bytesData, []byte("用户名或密码不正确，请重新输入！"))
	if ind != -1 {
		z.errs = append(z.errs, ErrUsernameOrPasswordWrong)
		return ErrUsernameOrPasswordWrong
	}
	return nil
}

type KcItem struct {
	Bfzcj              int    `json:"bfzcj"`     // 百分制 成绩
	BhId               string `json:"bh_id"`     // 编号id
	Bj                 string `json:"bj"`        // 班级信息
	Cj                 string `json:"cj"`        // 成绩 字符串
	Cjsfzf             string `json:"cjsfzf"`    //
	Date               string `json:"date"`      // 大写日期
	DateDigit          string `json:"dateDigit"` // 数字日期
	DateDigitSeparator string `json:"dateDigitSeparator"`
	Day                string `json:"day"`   // 几天是几号 23号
	Jd                 string `json:"jd"`    // 绩点
	JgId               string `json:"jg_id"` //
	Jgmc               string `json:"jgmc"`  // 机构信息
	Kcbj               string `json:"kcbj"`  // 课程背景 主修
	Kch                string `json:"kch"`   // 课程号
	KchId              string `json:"kch_id"` // 课程号id
	Kclbmc             string `json:"kclbmc"` // ? 课程列表名称
	Kcmc               string `json:"kcmc"` // 课程名称
	Kcxzdm             string `json:"kcxzdm"`  // ?? 课程行政代码
	Kcxzmc             string `json:"kcxzmc"`  // 必修课
	Key                string `json:"key"`     // 课程key
	Kkbmmc             string `json:"kkbmmc"`  // 课程开设的学院
	Ksxz               string `json:"ksxz"`    // 考试类型
	Month              string `json:"month"`   // 月份
	NjmcId             string `json:"njmc_id"` //
	Njmc               string `json:"njmc"`
	Xf                 string `json:"xf"` // 学分

	Xb     string `json:"xb"`  // 性别
	Xbm    string `json:"xbm"` // 性别码
	Xh     string `json:"xh"`  // 学号
	XhId   string `json:"xh_id"`
	Xm     string `json:"xm"`    // 姓名
	Xnm    string `json:"xnm"`   // 2013
	Xnmmc  string `json:"xnmmc"` // 学期
	Xqm    string `json:"xqm"`
	Xqmmc  string `json:"xqmmc"`
	Xsbjmc string `json:"xsbjmc"` // 学校
	Xslb   string `json:"xslb"`   // 学生类别
	Year   string `json:"year"`   //
	ZyhId  string `json:"zyh_id"` // 专业id
	Zymc   string `json:"zymc"`   // 专业名称
}

type KcResult struct {
	CurrentPage   int      `json:"currentPage"`
	CurrentResult int      `json:"currentResult"`
	Limit         int      `json:"limit"`
	Offset        int      `json:"offset"`
	PageNo        int      `json:"pageNo"`
	PageSize      int      `json:"pageSize"`
	ShowCount     int      `json:"showCount"`
	SortOrder     string   `json:"sortOrder"`
	TotalCount    int      `json:"totalCount"`
	TotalPage     int      `json:"totalPage"`
	TotalResult   int      `json:"totalResult"`
	Items         []KcItem `json:"items"`
}

// TODO: 添加过滤器 筛选数据
func (z *ZSpider) GetKcs() (*KcResult, error) {
	if len(z.errs) != 0 {
		return nil, z.errs[0]
	}
	data := url.Values{
		"xh_id":                  {z.number},
		"xnm":                    {""},
		"xqm":                    {""},
		"_search":                {fmt.Sprintf("%v", false)},
		"nd":                     {fmt.Sprintf("%d", time.Now().UnixNano()/1e6)},
		"time":                   {"0"},
		"queryModel.showCount":   {"5000"},
		"queryModel.currentPage": {"1"},
		"queryModel.sortName":    {""},
		"queryModel.sortOrder":   {"asc"},
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(kcUrl, z.number), strings.NewReader(data.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := z.client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	bData, _ := ioutil.ReadAll(resp.Body)

	var results KcResult
	json.Unmarshal(bData, &results)
	return &results, nil
}
