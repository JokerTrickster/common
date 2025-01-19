package error

import (
	"context"
	"encoding/json"
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
	Trace   string `json:"trace,omitempty"`
	From    string `json:"from,omitempty"`
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

type CustomError struct {
	ErrType string `json:"errType"`
	Trace   string `json:"trace"`
	Msg     string `json:"msg"`
	From    string `json:"from"`
}

func (e CustomError) Error() string {
	// 에러 메시지를 JSON 문자열로 변환
	errorJSON, _ := json.Marshal(e)
	return string(errorJSON)
}

func CreateError(ctx context.Context, errType string, trace string, msg string, from string) error {
	return CustomError{
		ErrType: errType,
		Trace:   trace,
		Msg:     msg,
		From:    from,
	}
}

// Trace captures the function name and line number where the error occurred
func Trace() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	_, line := runtime.FuncForPC(pc).FileLine(pc)
	return fmt.Sprintf("%s.L%d", funcName, line)
}

// GenerateHTTPErrorResponse generates a standard HTTP error response
func GenerateHTTPErrorResponse(err error) (int, ResError) {
	// Parse the error into an Err struct
	parsedErr := ParseError(err.Error())

	// Map the parsed error to an HTTP response
	resError := ResError{
		ErrType: parsedErr.ErrType,
		Msg:     parsedErr.Msg,
		Trace:   parsedErr.Trace,
		From:    parsedErr.From,
	}
	return parsedErr.HttpCode, resError
}

// GenerateCustomErrorResponse allows creating custom errors without using ParseError
func GenerateCustomErrorResponse(httpCode int, errType, msg string) (int, ResError) {
	resError := ResError{
		ErrType: errType,
		Msg:     msg,
	}
	return httpCode, resError
}
