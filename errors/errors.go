/*
自定义的错误处理，基于code生成
*/
package errors

import (
	"fmt"
)

type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

// Error return the externally-safe error message.
func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

// Cause return the cause of the withCode error.
func (w *withCode) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withCode) Unwrap() error { return w.cause }

// 根据code生成error
func WithCode(code int, format string, args any) error {
	return &withCode{
		err:   fmt.Errorf(format, args),
		code:  code,
		stack: callers(),
	}
}

// 构建错误链
func WrapC(err error, code int, format string, args any) error {
	if err == nil {
		return nil
	}
	return &withCode{
		err:   fmt.Errorf(format, args),
		code:  code,
		cause: err,
		stack: callers(),
	}
}
