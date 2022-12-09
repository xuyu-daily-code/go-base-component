package log

import (
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type JsonFormatter struct {
	IgnoreBasicFields bool
}

func (f *JsonFormatter) Format(entry *Entry) error {
	if !f.IgnoreBasicFields {
		entry.Map["level"] = LevelNameMapping[entry.Level]
		entry.Map["time"] = entry.Time.Format(time.RFC3339)
		if entry.File != "" {
			entry.Map["file"] = entry.File + ":" + strconv.Itoa(entry.Line)
			entry.Map["func"] = entry.Func
		}

		switch entry.Format {
		case FmtEmptySeparate:
			entry.Map["message"] = fmt.Sprint(entry.Args...)
		default:
			entry.Map["message"] = fmt.Sprintf(entry.Format, entry.Args...)
		}
		return jsoniter.NewEncoder(entry.Buffer).Encode(entry.Map)
	}

	switch entry.Format {
	case FmtEmptySeparate:
		for _, arg := range entry.Args {
			if err := jsoniter.NewEncoder(entry.Buffer).Encode(arg); err != nil {
				return err
			}
		}
	default:
		entry.Buffer.WriteString(fmt.Sprintf(entry.Format, entry.Args...))
	}
	return nil
}
