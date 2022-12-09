package log

import (
	"io"
	"os"
)

type Level uint8

const (
	FmtEmptySeparate = ""
)

// 日志等级
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// 日志等级和中文描述映射
var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

// 日志选项定义
type options struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
}

type Option func(*options)

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = os.Stdout
	}

	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	return o
}

// 设置输出位置
func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

// 设置输出级别
func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithStdLevel(stdLevel Level) Option {
	return func(o *options) {
		o.stdLevel = stdLevel
	}
}

// 设置输出格式
func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

// 设置是否打印文件名和行号
func WithDisableCaller(disableCaller bool) Option {
	return func(o *options) {
		o.disableCaller = disableCaller
	}
}
