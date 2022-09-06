package log

import (
	"go.uber.org/zap"
	"sync"
)

type Logger struct {
	logger *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
}

var (
	std = New(NewOptions())
	mu  = sync.Mutex{}
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = New(opts)

	// 需要将导出的函数重新赋值为新创建的logger的函数
	Info = std.Info
	Infof = std.Infof
	Warn = std.Warn
	Warnf = std.Warnf
	Error = std.Error
	Errorf = std.Errorf
	DPanic = std.DPanic
	DPanicf = std.DPanicf
	Panic = std.Panic
	Panicf = std.Panicf
	Fatal = std.Fatal
	Fatalf = std.Fatalf
	Debug = std.Debug
	Debugf = std.Debugf
	Sync = std.Sync
}

func New(opts *Options) *Logger {
	return opts.Build()
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Sugar().Warnf(format, v)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Sugar().Errorf(format, v)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

func (l *Logger) DPanicf(format string, v ...interface{}) {
	l.logger.Sugar().DPanicf(format, v)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logger.Sugar().Panicf(format, v)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Sugar().Fatalf(format, v)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}
