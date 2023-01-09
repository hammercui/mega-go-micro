/*
@Desc :  自定义recovery中间件
@Version : 1.0.0
@Time : 2019/4/17 15:17
@Author : hammercui
@File : gorecovery
@Company: Sdbean
*/
package gin

import (
	"bytes"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/hammercui/mega-go-micro/v2/log"
	"io/ioutil"
	"net/http"
	//"net/http/httputil"
	"runtime"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//if logger != nil {
				Stack(3)
				//httprequest, _ := httputil.DumpRequest(c.Request, false)
				httprequest := c.Request.URL.Path
				logStr := fmt.Sprintf("[panic cause]%s [%s]\n %s\n", err, timeFormat(time.Now()),
					httprequest)
				log.Logger().Error(logStr)
				//sentry.CaptureException(err.(error))
				//sentry.Flush(5 * time.Second)
				//}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

//提供给其他协程捕捉错误
func CatchPanic()  {
	if err := recover(); err != nil {
		Stack(3)
		log.Logger().Error(err)
		sentry.CaptureException(err.(error))
		sentry.Flush(2 * time.Second)
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	var panicLog string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		//fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		//log.Logger().Error(fmt.Sprintf("[panic stack] %s:%d (0x%x)",file, line, pc))
		panicLog = panicLog + fmt.Sprintf("[panic stack] %s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		//fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
		//log.Logger().Error(fmt.Sprintf("[panic stack] \t%s: %s\n",function(pc), source(lines, line)))
		panicLog = panicLog + fmt.Sprintf("[panic stack] \t%s: %s\n", function(pc), source(lines, line))
	}
	log.Logger().Error(panicLog)
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
