package errcode

import "fmt"

// Error 错误码结构
type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details,omitempty"`
}

// NewError 创建错误码
func NewError(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Msg)
}

// GetCode 获取错误码
func (e *Error) GetCode() int {
	return e.Code
}

// GetMsg 获取错误信息
func (e *Error) GetMsg() string {
	return e.Msg
}

// WithDetails 添加错误详情
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.Details = []string{}
	newError.Details = append(newError.Details, details...)
	return &newError
}
