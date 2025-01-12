package error

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

/*
	공통 에러 타입 및 유틸리티 함수를 정의합니다.
*/

// ResError represents the error format for frontend
type ResError struct {
	ErrType string `json:"errType,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// Err represents the error format for logging
type Err struct {
	HttpCode int    `json:"httpCode,omitempty"`
	ErrType  string `json:"errType,omitempty"`
	Msg      string `json:"msg,omitempty"`
	Trace    string `json:"trace,omitempty"`
	From     string `json:"from,omitempty"`
}

// ParseError parses the error string into an Err struct
func ParseError(data string) Err {
	slice := strings.Split(data, "|")
	return Err{
		HttpCode: ErrHttpCode[slice[0]],
		ErrType:  slice[0],
		Trace:    slice[1],
		Msg:      slice[2],
		From:     slice[3],
	}
}

// CreateError formats an error with additional context information
func CreateError(ctx context.Context, errType string, trace string, msg string, from string) error {
	return fmt.Errorf("%s|%s|%s|%s", errType, trace, msg, from)
}

// Trace captures the function name and line number where the error occurred
func Trace() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	return fmt.Sprintf("%s.L%d", funcName, line)
}
