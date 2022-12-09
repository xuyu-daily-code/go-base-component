package log

import (
	"fmt"
	"time"
)

type TextFormatter struct {
	// 是否打印基础信息
	IgnoreBasicFields bool
}

func (f *TextFormatter) Format(entry *Entry) error {
	if !f.IgnoreBasicFields {
		// 拼装日志等级和日志时间
		entry.Buffer.WriteString(fmt.Sprintf("%s %s", entry.Time.Format(time.RFC3339), LevelNameMapping[entry.Level]))
		if entry.File != "" {
			short := entry.File
			for i := len(entry.File) - 1; i > 0; i-- {
				if entry.File[i] == '/' {
					short = entry.File[i+1:]
					break
				}
			}
			entry.Buffer.WriteString(fmt.Sprintf(" %s:%d", short, entry.Line))
		}
		entry.Buffer.WriteString(" ")
	}

	switch entry.Format {
	case FmtEmptySeparate:
		entry.Buffer.WriteString(fmt.Sprint(entry.Args...))
	default:
		entry.Buffer.WriteString(fmt.Sprintf(entry.Format, entry.Args...))
	}
	entry.Buffer.WriteString("\n")

	return nil
}
