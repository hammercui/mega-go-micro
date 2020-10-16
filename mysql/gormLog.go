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
	"wfServerMicro/infra/conf"
	infraLog "wfServerMicro/infra/log"
)

func NewGormLog(env conf.AppEnv) *logger {
	//执行超过5ms报警
	slowThresHold := 7 * time.Millisecond
	if env == conf.AppEnv_prod {
		//生产超过1ms报警
		slowThresHold = 5 * time.Millisecond
	}
	return &logger{
		SlowThreshold: slowThresHold,
		Env:           env,
	}
}

/**
**	1 info日志相当于debug
**  2 测试环境显示(info,warn,error,)
**  3 生产显示(warn,error)
 */
type logger struct {
	SlowThreshold time.Duration
	Env           conf.AppEnv
}

// LogMode log mode
func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newlogger := *l
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	//l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	if l.Env != conf.AppEnv_prod {
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
	case l.Env == conf.AppEnv_local:
		sql, rows := fc()
		//	//l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		infraLog.Logger().Debugf("GormDebug |%fms |rows:%d |sql:%s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
	//}
}
