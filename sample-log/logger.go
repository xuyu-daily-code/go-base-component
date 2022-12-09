package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"unsafe"
)

// 默认全局logger
var std = New()

type logger struct {
	opt       *options   // logger 启动参数
	mu        sync.Mutex // 锁
	entryPool *sync.Pool //临时对象存储池
}

func New(pts ...Option) *logger {
	logger := &logger{opt: initOptions(pts...)}
	logger.entryPool = &sync.Pool{New: func() any { return entry(logger) }}
	return logger
}

// 获取全局logger
func StdLogger() *logger {
	return std
}

// 设置启动参数
func (logger *logger) SetOptions(opts ...Option) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	for _, opt := range opts {
		opt(logger.opt)
	}
}

func (logger *logger) Writer() io.Writer {
	return logger
}

// logger对象实现io.Write接口
func (logger *logger) Write(data []byte) (int, error) {
	logger.entry().write(logger.opt.stdLevel, FmtEmptySeparate, *(*string)(unsafe.Pointer(&data)))
	return 0, nil
}

func (logger *logger) entry() *Entry {
	return logger.entryPool.Get().(*Entry)
}

func Writer() io.Writer {
	return std
}

/*
非全局logger非格式化输出日志接口
*/
func (logger *logger) Debug(args ...any) {
	logger.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func (logger *logger) Info(args ...any) {
	logger.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func (logger *logger) Warn(args ...any) {
	logger.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func (logger *logger) Error(args ...any) {
	logger.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func (logger *logger) Panic(args ...any) {
	logger.entry().write(PanicLevel, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func (logger *logger) Fatal(args ...any) {
	logger.entry().write(FatalLevel, FmtEmptySeparate, args...)
	os.Exit(1)
}

/*
非全局logger格式化输出日志接口
*/
func (logger *logger) Debugf(format string, args ...any) {
	logger.entry().write(DebugLevel, format, args...)
}

func (logger *logger) Infof(format string, args ...any) {
	logger.entry().write(InfoLevel, format, args...)
}

func (logger *logger) Warnf(format string, args ...any) {
	logger.entry().write(WarnLevel, format, args...)
}

func (logger *logger) Errorf(format string, args ...any) {
	logger.entry().write(ErrorLevel, format, args...)
}

func (logger *logger) Panicf(format string, args ...any) {
	logger.entry().write(PanicLevel, format, args...)
	panic(fmt.Sprint(args...))
}

func (logger *logger) Fatalf(format string, args ...any) {
	logger.entry().write(FatalLevel, format, args...)
	os.Exit(1)
}

/*
全局logger非格式化输出日志接口
*/
func Debug(args ...any) {
	std.entry().write(DebugLevel, FmtEmptySeparate, args...)
}

func Info(args ...any) {
	std.entry().write(InfoLevel, FmtEmptySeparate, args...)
}

func Warn(args ...any) {
	std.entry().write(WarnLevel, FmtEmptySeparate, args...)
}

func Error(args ...any) {
	std.entry().write(ErrorLevel, FmtEmptySeparate, args...)
}

func Panic(args ...any) {
	std.entry().write(PanicLevel, FmtEmptySeparate, args...)
	panic(fmt.Sprint(args...))
}

func Fatal(args ...any) {
	std.entry().write(FatalLevel, FmtEmptySeparate, args...)
	os.Exit(1)
}

/*
全局logger格式化输出日志接口
*/
func Debugf(format string, args ...any) {
	std.entry().write(DebugLevel, format, args...)
}

func Infof(format string, args ...any) {
	std.entry().write(InfoLevel, format, args...)
}

func Warnf(format string, args ...any) {
	std.entry().write(WarnLevel, format, args...)
}

func Errorf(format string, args ...any) {
	std.entry().write(ErrorLevel, format, args...)
}

func Panicf(format string, args ...any) {
	std.entry().write(PanicLevel, format, args...)
	panic(fmt.Sprint(args...))
}

func Fatalf(format string, args ...any) {
	std.entry().write(FatalLevel, format, args...)
	os.Exit(1)
}
