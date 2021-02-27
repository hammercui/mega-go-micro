/**
 * Description:skyWalking中间件
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/2/27
 * Time: 13:36
 * Mail: hammercui@163.com
 *
 */
package middleware

import (
	"fmt"
	"github.com/hammercui/go2sky"
	"github.com/hammercui/go2sky/propagation"
	v3 "github.com/hammercui/go2sky/reporter/grpc/language-agent"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
	"time"
)

const (
	httpServerComponentID int32 = 49
)

type routeInfo struct {
	operationName string
}

type middleware struct {
	routeMap     map[string]map[string]routeInfo
	routeMapOnce sync.Once
}

//Middleware gin middleware return HandlerFunc  with tracing.
func SkyWalking(engine *gin.Engine, tracer *go2sky.Tracer) gin.HandlerFunc {
	if engine == nil || tracer == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	m := new(middleware)

	return func(c *gin.Context) {
		m.routeMapOnce.Do(func() {
			routes := engine.Routes()
			rm := make(map[string]map[string]routeInfo)
			for _, r := range routes {
				mm := rm[r.Method]
				if mm == nil {
					mm = make(map[string]routeInfo)
					rm[r.Method] = mm
				}
				mm[r.Handler] = routeInfo{
					operationName: fmt.Sprintf("/%s%s", r.Method, r.Path),
				}
			}
			m.routeMap = rm
		})
		var operationName string
		handlerName := c.HandlerName()
		if routeInfo, ok := m.routeMap[c.Request.Method][handlerName]; ok {
			operationName = routeInfo.operationName
		}
		if operationName == "" {
			operationName = c.Request.Method
		}
		span, ctx, err := tracer.CreateEntrySpan(c.Request.Context(), operationName, func() (string, error) {
			return c.Request.Header.Get(propagation.Header), nil
		})
		if err != nil {
			c.Next()
			return
		}
		span.SetComponent(httpServerComponentID)
		span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
		span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
		span.SetSpanLayer(v3.SpanLayer_Http)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			span.Error(time.Now(), c.Errors.String())
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
		span.End()
	}
}
