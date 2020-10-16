/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2020/9/5
 * Time: 18:05
 * Mail: hammercui@163.com
 *
 */
package consul
import "context"

func watchStale(ctx context.Context) bool {
	if ctx == nil {
		return true
	}

	stale, ok := ctx.Value(watchStaleKey{}).(bool)
	if !ok {
		return true
	}
	return stale
}