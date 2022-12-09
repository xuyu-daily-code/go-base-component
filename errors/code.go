package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder defaultCoder = defaultCoder{
		C:    1,
		HTTP: http.StatusInternalServerError,
		Ext:  "An internal server error occurred",
		Ref:  "",
	}
)

type Coder interface {
	// 返回http code
	HTTPStatus() int
	// 返回错误信息
	String() string
	// 返回详细的文档地址
	Reference() string
	// 返回错误码
	Code() int
}

// 默认code
type defaultCoder struct {
	C    int
	HTTP int
	Ext  string
	Ref  string
}

/*
defaultCoder 实现接口Coder
*/
func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}
	return coder.HTTP
}

func (coder defaultCoder) String() string {
	return coder.Ext
}

func (coder defaultCoder) Reference() string {
	return coder.Ref
}

func (coder defaultCoder) Code() int {
	return coder.C
}

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// 注册错误码和对象
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/xyedu/errors` as unknownCode error code")
	}
	codeMux.Lock()
	defer codeMux.Unlock()
	codes[coder.Code()] = coder
}

// 注册错误码和对象，重复抛出panic
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved by `github.com/xyedu/errors` as unknownCode error code")
	}
	codeMux.Lock()
	defer codeMux.Unlock()
	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}
	codes[coder.Code()] = coder
}

func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}
	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}
	return unknownCoder
}

// 判断code是否在错误链中
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}
		if v.cause != nil {
			return IsCode(v.cause, code)
		}
		return false
	}
	return false
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}
