// 自定义的错误输出格式化
package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// 包含所有错误信息
type formatInfo struct {
	code    int
	message string
	err     string
	stack   *stack
}

// 格式化规则
// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Verbs:
//
//	%s  - Returns the user-safe error string mapped to the error code or
//	  ┊   the error message if none is specified.
//	%v      Alias for %s
//
// Flags:
//
//	#      JSON formatted output, useful for logging
//	-      Output caller details, useful for troubleshooting
//	+      Output full error stack details, useful for debugging
//
// Examples:
//
//	%s:    error for internal read B
//	%v:    error for internal read B
//	%-v:   error for internal read B - #0 [/home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:12 (main.main)] (#100102) Internal Server Error
//	%+v:   error for internal read B - #0 [/home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:12 (main.main)] (#100102) Internal Server Error; error for internal read A - #1 [/home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:35 (main.newErrorB)] (#100104) Validation failed
//	%#v:   [{"error":"error for internal read B"}]
//	%#-v:  [{"caller":"#0 /home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:12 (main.main)","error":"error for internal read B","message":"(#100102) Internal Server Error"}]
//	%#+v:  [{"caller":"#0 /home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:12 (main.main)","error":"error for internal read B","message":"(#100102) Internal Server Error"},{"caller":"#1 /home/lk/workspace/golang/src/github.com/marmotedu/iam/main.go:35 (main.newErrorB)","error":"error for internal read A","message":"(#100104) Validation failed"}]
func (w *withCode) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		jsonData := []map[string]any{}
		var (
			flagDetail bool
			flagTrace  bool
			modelJson  bool
		)

		if state.Flag('#') {
			modelJson = true
		}

		if state.Flag('-') {
			flagDetail = true
		}

		if state.Flag('+') {
			flagTrace = true
		}

		sep := ""
		errs := list(w)
		length := len(errs)
		for k, e := range errs {
			finfo := buildFormatInfo(e)
			jsonData, str = format(length-k-1, jsonData, str, finfo, sep, flagDetail, flagTrace, modelJson)
			sep = ";"
			if !flagTrace {
				break
			}
		}
		if modelJson {
			bytes, _ := json.Marshal(jsonData)
			str.Write(bytes)
		}
		fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
	default:
		finfo := buildFormatInfo(w)
		fmt.Fprintf(state, finfo.message)
	}
}

// 将formatInfo对象格式化到指定的输出格式
func format(k int, jsonData []map[string]any, str *bytes.Buffer, finfo *formatInfo,
	sep string, flagDetail, flagTrace, modelJson bool) ([]map[string]any, *bytes.Buffer) {
	if modelJson {
		data := map[string]any{}
		if flagDetail || flagTrace {
			data = map[string]any{
				"message": finfo.message,
				"code":    finfo.code,
				"error":   finfo.err,
			}
			// 错误栈序号
			caller := fmt.Sprintf("#%d", k)
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				caller = fmt.Sprintf("%s %s:%d (%s)",
					caller,
					f.file(),
					f.line(),
					f.name(),
				)
			}
			data["caller"] = caller
		} else {
			data["error"] = finfo.message
		}
		jsonData = append(jsonData, data)
	} else {
		if flagDetail || flagTrace {
			if finfo.stack != nil {
				f := Frame((*finfo.stack)[0])
				fmt.Fprintf(str, "%s%s - #%d [%s:%d (%s)] (%d) %s",
					sep,
					finfo.err,
					k,
					f.file(),
					f.line(),
					f.name(),
					finfo.code,
					finfo.message,
				)
			} else {
				fmt.Fprintf(str, "%s%s - #%d %s", sep, finfo.err, k, finfo.message)
			}
		} else {
			fmt.Fprintf(str, finfo.message)
		}
	}

	return jsonData, str
}

// 将error由链式转换为数组
func list(e error) []error {
	ret := []error{}
	if e != nil {
		if w, ok := e.(interface{ Unwrap() error }); ok {
			ret = append(ret, e)
			ret = append(ret, list(w.Unwrap())...)
		} else {
			ret = append(ret, e)
		}
	}
	return ret
}

// 由error对象组装formatInfo对象
func buildFormatInfo(err error) *formatInfo {
	var finfo *formatInfo
	switch err := err.(type) {
	case *withCode:
		coder, ok := codes[err.code]
		if !ok {
			coder = unknownCoder
		}
		extMsg := coder.String()
		if extMsg == "" {
			extMsg = err.err.Error()
		}
		finfo = &formatInfo{
			code:    coder.Code(),
			message: extMsg,
			err:     err.err.Error(),
			stack:   err.stack,
		}
	default:
		finfo = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
		}
	}
	return finfo
}
