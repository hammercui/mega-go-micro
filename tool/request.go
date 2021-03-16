/*
@Desc : 2020/8/21 14:15
@Version : 1.0.0
@Time : 2020/8/21 14:15
@Author : hammercui
@File : request
@Company: Sdbean
*/
package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hammercui/mega-go-micro/conf"
	"github.com/hammercui/mega-go-micro/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type RequestOptions struct {
	//超时时间
	Timeout time.Duration
	//ginContext
	ginCtx *gin.Context
	//ctx

}

//默认配置
var DefaultRequestOptions = &RequestOptions{
	//默认超时10秒
	Timeout: 10 * time.Second,
}

//生成http请求签名
func genReqSign() string {
	timeStamp := time.Now().UnixNano() / 1e6
	nodeIdStr := "1"
	if conf.GetConf() != nil && conf.GetConf().AppConf != nil {
		nodeIdStr = conf.GetConf().AppConf.NodeId
	}
	nodeId, _ := strconv.Atoi(nodeIdStr)
	offset := int64(nodeId) * int64(1000000000000)
	newTimeStamp := timeStamp + offset
	return fmt.Sprintf("%d", newTimeStamp)
}

//使用默认配置的json请求
func PostJson(url string, v interface{}, out interface{}) error {
	return PostJsonWithOpt(url, v, out, DefaultRequestOptions)
}

//自定义配置的json请求
func PostJsonWithOpt(url string, v interface{}, out interface{}, opts *RequestOptions) error {
	reqSign := genReqSign()
	log.Logger().Infof("[%s]http request-->: url[%s]", reqSign, url)
	log.Logger().Infof("[%s]http request->:timeout:%v", reqSign, opts.Timeout)
	client := &http.Client{
		Timeout: opts.Timeout,
	}
	//处理span探针

	var bytesData []byte
	//v是字符串
	if vstr, ok := v.(string); ok {
		bytesData = []byte(vstr)
	} else {
		if vbyte, err := json.Marshal(v); err == nil {
			bytesData = vbyte
		} else {
			log.Logger().Errorf("[%s]json Marshal err:%+v", reqSign, err)
			return err
		}
	}
	log.Logger().Infof("[%s]http request->:%s", reqSign, string(bytesData))
	//logger.Debug(fmt.Sprintf("http request-->url: %s",url))
	//logger.Debug(fmt.Sprintf("http request-->data: %s",string(bytesData)))
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Logger().Errorf("[%s]http request err:%+v", reqSign, err)
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(request)
	if err != nil {
		log.Logger().Errorf("[%s]http request do err:%+v", reqSign, err)
		return err
	}
	//http !=200
	if resp.StatusCode != http.StatusOK {
		log.Logger().Errorf("[%s]http request do err: statusCode=%d", reqSign, resp.StatusCode)
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		log.Logger().Errorf("[%s]resp.Body close err:%+v", reqSign, err)
		return err
	}
	if err != nil {
		log.Logger().Errorf("[%s]respBytes err:%+v", reqSign, err)
		return err
	}
	log.Logger().Infof("[%s]http response<--: %s", reqSign, string(respBytes))
	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	if out == nil {
		return nil
	}
	err = json.Unmarshal(respBytes, out)
	if err != nil {
		log.Logger().Errorf("[%s]http json Unmarshal err", reqSign, err)
		return err
	}
	return nil
}

//请求表单
func PostForm(urlStr string, data url.Values, out interface{}) error {
	//默认值
	timeout := 10 * time.Second
	log.Logger().Infof("http request-->: url[%s]", urlStr)
	log.Logger().Info("http request->:timeout", timeout)
	client := &http.Client{
		Timeout: timeout,
	}
	reqBody := strings.NewReader(data.Encode())
	r, _ := http.NewRequest("POST", urlStr, reqBody) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	log.Logger().Info("http request->:body", data)

	resp, err := client.Do(r)
	if err != nil {
		log.Logger().Error("request do err: ", err)
		return err
	}
	defer resp.Body.Close()
	//http !=200
	if resp.StatusCode != http.StatusOK {
		log.Logger().Error("request do err: statusCode=%d", resp.StatusCode)
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err := resp.Body.Close(); err != nil {
		log.Logger().Error("resp.Body close err:%+v", err)
		return err
	}
	if err != nil {
		log.Logger().Error("respBytes err:%+v", err)
		return err
	}
	log.Logger().Infof("http response<--: %s", string(respBytes))
	//byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	if out == nil {
		return nil
	}
	err = json.Unmarshal(respBytes, out)
	if err != nil {
		log.Logger().Error("json Unmarshal err", err)
		return err
	}
	return nil

}
