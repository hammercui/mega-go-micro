/**
 * Description: http中间件转换工具
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/3/1
 * Time: 17:38
 * Mail: hammercui@163.com
 *
 */
package gin

import (
	"context"
	"github.com/gin-gonic/gin"
)

//context转换为*gin.Context
func ConvertGinContext(ctx context.Context) *gin.Context {
	if p, ok := ctx.(*gin.Context); ok {
		return p
	}
	return nil
}
