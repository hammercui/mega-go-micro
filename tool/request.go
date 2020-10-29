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
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"github.com/hammercui/mega-go-micro/log"
)

//从body获得json v2版本
func ReadPostJsonV2(c *gin.Context) []byte {
	body, _ := ioutil.ReadAll(c.Request.Body)
	return body
}

func PostJson(url string, v interface{},out interface{}) error  {
	//默认值
	timeout := 10 * time.Second

	//大于配置毫秒，按传值计算
	//if (millSec > int(conf.Config.Server.PhpTimeOut)) {
	//	timeout = time.Duration(millSec) * time.Millisecond
	//}
	log.Logger().Infof("http request-->: url[%s]", url)
	log.Logger().Info("http request->:timeout", timeout)
	client := &http.Client{
		Timeout: timeout,
	}

	var bytesData []byte
	//v是字符串
	if vstr, ok := v.(string); ok {
		bytesData = []byte(vstr)
	} else {
		if vbyte, err := json.Marshal(v); err == nil {
			bytesData = vbyte
		} else {
			log.Logger().Error("json Marshal err", err)
			return err
		}
	}
	log.Logger().Info("http request->:body", string(bytesData))
	//logger.Debug(fmt.Sprintf("http request-->url: %s",url))
	//logger.Debug(fmt.Sprintf("http request-->data: %s",string(bytesData)))
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Logger().Error("request err:%+v", err)
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(request)
	if err != nil {
		log.Logger().Error("request do err:%+v", err)
		return err
	}
	//http !=200
	if (resp.StatusCode != http.StatusOK) {
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
	if(out == nil){
		return nil
	}
	err = json.Unmarshal(respBytes, out)
	if err != nil {
		log.Logger().Error("json Unmarshal err", err)
		return err;
	}
	return nil

}

//请求表单
func PostForm(urlStr string, data url.Values,out interface{}) error  {
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
	if (resp.StatusCode != http.StatusOK) {
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
	if(out == nil){
		return nil
	}
	err = json.Unmarshal(respBytes, out)
	if err != nil {
		log.Logger().Error("json Unmarshal err", err)
		return err;
	}
	return nil

}