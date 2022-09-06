package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"

	defaultLogPath     = "./logs"
	defaultLogFilename = "zap.log"
)

// Options 日志配置选项
type Options struct {
	// 日志级别，优先级从低到高依次为：Debug, Info, Warn, Error, Dpanic, Panic, Fatal
	Level string
	// 支持的日志输出格式，支持 json 和 console (就是 text 格式)
	Format string
	// 是否开启颜色输出
	EnableColor bool
	// 日志文件路径
	Path string
	// 日志文件名称
	Filename string
	// 单个文件最大大小，单位：MB
	MaxSize int
	// 旧日志保存的最大天数
	MaxAge int
	// 旧日志保存的最大个数
	MaxBackups int
	// 对backup日志是否压缩
	Compress bool
	// 是否开启 caller，如果开启会在日志中显示调用日志所在的文件、函数和行号
	EnableCaller bool
	// 是否在 Panic 及以上级别禁止打印堆栈信息
	EnableStacktrace bool
	// 是否输出到控制台
	EnableStdout bool
}

// NewOptions 通过默认参数创建 Options
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.DebugLevel.String(),
		EnableCaller:     true,
		EnableStacktrace: true,
		Format:           consoleFormat,
		EnableColor:      true,
		EnableStdout:     true,
		Path:             defaultLogPath,
		Filename:         defaultLogFilename,
		MaxSize:          100,
		MaxAge:           10,
		MaxBackups:       10,
		Compress:         false,
	}
}

func (o *Options) validate() {
	if o.Path == "" || o.Filename == "" {
		panic("请指定日志文件路径或名称")
	}
}

// Build 创建全局 zap logger
func (o *Options) Build() *Logger {
	o.validate()

	// zapLevel 是 int8 类型，具有零值
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoder := o.buildEncoder()
	writer, err := o.buildWriter()
	if err != nil {
		panic(err)
	}

	zapLogger := zap.New(
		zapcore.NewCore(encoder, writer, zapLevel),
		o.buildZapOptions()...,
	)

	return &Logger{
		logger: zapLogger,
	}
}

func (o *Options) buildEncoder() zapcore.Encoder {
	encodeLevel := zapcore.CapitalLevelEncoder
	// 文本格式并且输出到标准输出才开启颜色支持
	if o.Format == consoleFormat && o.EnableColor && o.EnableStdout {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = encodeLevel
	encoderConfig.EncodeDuration = milliSecondsDurationEncoder

	if o.Format == jsonFormat {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (o *Options) buildWriter() (zapcore.WriteSyncer, error) {
	if err := createLogPath(o.Path); err != nil {
		return nil, err
	}

	// 日志轮转/分割
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(o.Path, o.Filename),
		MaxSize:    o.MaxSize,
		MaxAge:     o.MaxAge,
		MaxBackups: o.MaxBackups,
		LocalTime:  false,
		Compress:   o.Compress,
	}

	if o.EnableStdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberjackLogger), zapcore.AddSync(os.Stdout)), nil
	}

	return zapcore.AddSync(lumberjackLogger), nil
}

// 根据配置构建 zap.option
func (o *Options) buildZapOptions() []zap.Option {
	var opts = []zap.Option{zap.AddCallerSkip(1)}

	if o.EnableCaller {
		opts = append(opts, zap.AddCaller())
	}

	if o.EnableStacktrace {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}

	return opts
}

// 创建日志文件目录
func createLogPath(path string) error {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}
