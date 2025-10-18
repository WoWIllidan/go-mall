package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	ctx     context.Context
	traceId string
	spanId  string
	pSpanId string
	_logger *zap.Logger
}

func (l *logger) Debug(msg string, kv ...interface{}) {
	l.log(zapcore.DebugLevel, msg, kv...)
}

func (l *logger) Info(msg string, kv ...interface{}) {
	l.log(zapcore.InfoLevel, msg, kv...)
}

func (l *logger) Warn(msg string, kv ...interface{}) {
	l.log(zapcore.WarnLevel, msg, kv...)
}

func (l *logger) Error(msg string, kv ...interface{}) {
	l.log(zapcore.ErrorLevel, msg, kv...)
}

// kv 应该是成对的数据, 类似: name,张三,age,10,...
func (l *logger) log(lvl zapcore.Level, msg string, kv ...interface{}) {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志行信息中增加追踪参数
	kv = append(kv, "traceid", l.traceId, "spanid", l.spanId, "pspanid", l.pSpanId)
	// 增加日志调用者信息, 方便查日志时定位程序位置
	funcName, file, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func (l *logger) getLoggerCallerInfo() (funcName, file string, line int) {

	pc, file, line, ok := runtime.Caller(3) // 回溯拿调用日志方法的业务函数的信息
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}

func New(ctx context.Context) *logger {
	var traceId, spanId, pSpanId string
	if ctx.Value("traceid") != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if ctx.Value("psapnid") != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}

	return &logger{
		ctx:     ctx,
		traceId: traceId,
		spanId:  spanId,
		pSpanId: pSpanId,
		_logger: _logger,
	}
}
