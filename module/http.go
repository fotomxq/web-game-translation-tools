package module

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpType struct {
}

var (
	Http HttpType
)

//发送请求封装
// 配置下不需要提供密钥，直接就能获取到配置信息
//param getURL string 获取URL
//return error 是否成功
func (this *HttpType) HttpGet(getURL string) ([]byte, error) {
	//初始化参数
	var err error
	res := []byte{}
	//请求数据
	resp, err := http.DefaultClient.Get(getURL)
	if err != nil {
		return res, errors.New("resp is error , " + err.Error())
	}
	defer resp.Body.Close()
	//确定反馈信息必须是200，否则代表失败
	status := resp.StatusCode
	if status != 200 {
		return res, errors.New("status is not 200.")
	}
	//解析内容并反馈
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, errors.New("cannot read res.body , error : " + err.Error())
	}
	//反馈数据
	return res, nil
}

//发送请求封装
// 配置下不需要提供密钥，直接就能获取到配置信息
//param getURL string 获取URL
//param params url.Values 参数集合
//return error 是否成功
func (this *HttpType) HttpPost(getURL string, params url.Values) ([]byte, error) {
	//初始化参数
	var err error
	res := []byte{}
	//请求数据
	resp, err := http.DefaultClient.PostForm(getURL, params)
	if err != nil {
		return res, errors.New("resp is error , " + err.Error())
	}
	defer resp.Body.Close()
	//确定反馈信息必须是200，否则代表失败
	status := resp.StatusCode
	if status != 200 {
		return res, errors.New("status is not 200.")
	}
	//解析内容并反馈
	res, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, errors.New("cannot read res.body , error : " + err.Error())
	}
	//反馈数据
	return res, nil
}
