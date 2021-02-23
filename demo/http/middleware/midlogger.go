/*
@Desc : 日志中间件，打印路由信息
@Version : 1.0.0
@Time : 2019/4/17 14:10
@Author : hammercui
@File : midlogger
@Company: Sdbean
*/
package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
	"github.com/hammercui/mega-go-micro/log"
)

var (
	green        = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white        = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow       = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red          = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue         = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta      = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan         = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset        = string([]byte{27, 91, 48, 109})
	disableColor = false
)

var skipSwagger = regexp.MustCompile(`/swagger/*`)

//自定义日志
func Logger() gin.HandlerFunc {


	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if(path == "/"){
			c.Next()
			return
		}
		if strings.Contains(path, "/actuator") {
			c.Next()
			return
		}
		//copy入参
		reqBodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
		// Process request
		c.Next()
		//跳过swagger
		if skipSwagger.MatchString(path) {
			return
		}
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		//comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		logstr := fmt.Sprintf("%-7s | %s |%3d |%13v |%s |%s",
			method,
			path,
			statusCode,
			latency,
			clientIP,
			string(reqBodyBytes))
		//log := fmt.Sprintf("%s %-7s %s| %s |%s %3d %s|%13v|%15s %s",
		//	methodColor, method,resetColor,
		//	path,
		//	statusColor, statusCode,resetColor,
		//	latency,clientIP,comment)
		//fmt.Println("log内容",log)
		//logger.Debug(log)
		//fmt.Println(log)
		switch statusCode {
		case 200:
			log.Logger().Info(logstr)
		case 400:
			log.Logger().Warn(logstr)
		case 404:
			log.Logger().Error(logstr)
		default:
			log.Logger().Error(logstr)
		}
	}
}

//func colorForStatus(code int) string {
//	switch {
//	case code >= http.StatusOK && code < http.StatusMultipleChoices:
//		return green
//	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
//		return white
//	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
//		return yellow
//	default:
//		return red
//	}
//}
