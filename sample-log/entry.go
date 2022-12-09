package log

import (
	"bytes"
	"runtime"
	"strings"
	"time"
)

/*
保存所有的日志信息
日志配置和日志内容
*/
type Entry struct {
	// logger对象
	logger *logger
	// 缓存buffer
	Buffer *bytes.Buffer
	Map    map[string]any
	// 日志打印等级
	Level Level
	// 日志时间
	Time time.Time
	//
	File   string
	Line   int
	Func   string
	Format string
	Args   []any
}

func entry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		Buffer: new(bytes.Buffer),
		Map:    make(map[string]any, 5),
	}
}

func (e *Entry) write(level Level, format string, args ...any) {
	// 配置的日志等级大于当前打印等级，就不打印直接返回
	if e.logger.opt.level > level {
		return
	}
	// 日志时间为当前时间
	e.Time = time.Now()
	e.Level = level
	e.Format = format
	e.Args = args
	// 打印行数和方法
	if !e.logger.opt.disableCaller {
		if pc, file, line, ok := runtime.Caller(2); !ok {
			e.File = "???"
			e.Func = "???"
		} else {
			e.File, e.Line, e.Func = file, line, runtime.FuncForPC(pc).Name()
			e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}
	e.format()
	e.writer()
	e.release()
}

// 格式化
func (e *Entry) format() {
	_ = e.logger.opt.formatter.Format(e)
}

// 日志写入
func (e *Entry) writer() {
	e.logger.mu.Lock()
	_, _ = e.logger.opt.output.Write(e.Buffer.Bytes())
	e.logger.mu.Unlock()
}

// 清空缓存池和对象
func (e *Entry) release() {
	e.Args, e.Line, e.File, e.Format, e.Func = nil, 0, "", "", ""
	e.Buffer.Reset()
	e.logger.entryPool.Put(e)
}
