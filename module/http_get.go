package module

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//网络套件

type HttpGetType struct {
	//Header头预定义信息
	HTTPGetUserAgents []string
}

var (
	HttpGet HttpGetType
)

func (self *HttpGetType) SetConfig() {
	self.HTTPGetUserAgents = []string{
		"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
		"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
		"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
		"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	}
}

//查询HTML中存在多少个字符串
//param sendURL string URL地址
//param findStr string 要搜索的字符串
//return int 存在多少个，-1为失败，0没找到
func (self *HttpGetType) GetGoqueryStrInHtml(doc *goquery.Document, findStr string) int {
	html, err := doc.Html()
	if err != nil {
		return -1
	}
	return strings.Count(html, findStr)
}

//通过GET或POST获取URL数据
//param getURL string 获取URL地址
//param params url.Values 表单参数，只有post给定，留空则认定为get模式
//param proxyIP string 代理IP地址，如果留空跳过
//param isSetHeader bool 是否加入头信息加密，建议爬虫使用
//return []byte,error 数据，错误
func (self *HttpGetType) GetData(getURL string, params url.Values, proxyIP string, isSetHeader bool) ([]byte, error) {
	resp, err := self.GetResp(getURL, params, proxyIP, isSetHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		robots, err := ioutil.ReadAll(resp.Body)
		//resp.Body.Close()
		if err != nil {
			return []byte{}, err
		}
		return robots, nil
	} else {
		return []byte{}, err
	}
	return []byte{}, err
}

//获取RESP信息源
// 注意自行增加关闭机制
//param getURL string 获取URL地址
//param params url.Values 表单参数，只有post给定，留空则认定为get模式
//param proxyIP string 代理IP地址，如果留空跳过
//param isSetHeader bool 是否加入头信息加密，建议爬虫使用
//return *http.Response,error 数据，错误
func (self *HttpGetType) GetResp(getURL string, params url.Values, proxyIP string, isSetHeader bool) (*http.Response, error) {
	//初始化参数
	var resp *http.Response
	var urlx *url.URL
	var err error
	//设定代理
	client := &http.Client{}
	if proxyIP != "" {
		urlproxy, _ := urlx.Parse(proxyIP)
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
	}
	//设定反馈头
	var req *http.Request
	if params == nil {
		req, err = http.NewRequest(http.MethodGet, getURL, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, getURL, nil)
		req.Form = params
	}
	if err != nil {
		return nil, err
	}
	//如果需要对头信息加密，则进行加密处理
	if isSetHeader == true {
		req.Header.Add("User-Agent", self.GetUserAgentRand())
	}
	//执行URL获取
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	//定位结果
	status := resp.StatusCode
	if status != 200 {
		return resp, err
	}
	//defer resp.Body.Close()
	return resp, err
}

//通过goquery获取HTML
//param getURL string 获取URL地址
//param params url.Values POST参数，如果留空则GET
//param proxyIP string 代理IP地址，如果留空跳过
//param isSetHeader bool 是否加入头信息加密，建议爬虫使用
//return *goquery.Document , error 文档操作句柄，是否成功
func (self *HttpGetType) GetGoquery(getURL string, params url.Values, proxyIP string, isSetHeader bool) (*goquery.Document, error) {
	resp, err := self.GetResp(getURL, params, proxyIP, isSetHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return &goquery.Document{}, errors.New("GetGoquery Page Status is not 200.")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	return doc, err
}

//URL编码工具
func GetURLEncode(sendURL string) string {
	return url.QueryEscape(sendURL)
}

//随机获取一个UserAgent头信息
//return string 伪造头信息
func (self *HttpGetType) GetUserAgentRand() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return self.HTTPGetUserAgents[r.Intn(len(self.HTTPGetUserAgents))]
}
