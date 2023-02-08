package mysql

/**
 * Description:自定义gorm日志
 * version 1.0.0
 * Created by ${PRODUCT_NAME}.
 * Company sdbean
 * Author: hammercui
 * Date: ${DATE}
 * Time: ${TIME}
 * Mail: hammercui@163.com
 *
 */
import (
	"context"
	gormLogger "gorm.io/gorm/logger"
	"time"
	"github.com/hammercui/mega-go-micro/v2/conf"
	infraLog "github.com/hammercui/mega-go-micro/v2/log"
)

func NewGormLog(c *conf.MysqlConf) *logger {
	//默认执行超过10ms报警
	slowThreshold := 10 * time.Millisecond

	//读取配置
	if c.WarnThreshold > 0 {
		slowThreshold = time.Duration(c.WarnThreshold) * time.Millisecond
	}

	return &logger{
		SlowThreshold: slowThreshold,
		//Env:           configs.AppConf.Env,
		DebugInfo: c.DebugInfo,
	}
}

/**
**	1 info日志相当于debug
**  2 测试环境显示(info,warn,error,)
**  3 生产显示(warn,error)
 */
type logger struct {
	SlowThreshold time.Duration
	//Env           conf.AppEnv
	DebugInfo bool
}

// LogMode log mode
func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newlogger := *l
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	if l.DebugInfo {
		infraLog.Logger().Debug("GormDebug ", msg, data)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	//l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	infraLog.Logger().Warn("GormWarn ", msg, data)
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	//l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	infraLog.Logger().Error("GormError ", msg, data)
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	//if l.LogLevel > 0 {
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		//l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		infraLog.Logger().Errorf("GormError |err:%v |%fms |rows:%d |sql:%s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		sql, rows := fc()
		//l.Printf(l.traceWarnStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		infraLog.Logger().Warnf("GormWarn |%fms |rows:%d |sql:%s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.DebugInfo:
		sql, rows := fc()
		//	//l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		infraLog.Logger().Debugf("GormDebug |%fms |rows:%d |sql:%s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
	//}
}
