package validator

/*
	구조체 및 요청/응답 검증 로직
*/

import (
	"context"

	_errors "github.com/JokerTrickster/common/error" // 공통 에러 패키지 사용
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator instance
var val = validator.New()

// ValidateStruct validates a given struct
func ValidateStruct(class interface{}) error {
	if err := val.Struct(class); err != nil {
		return _errors.CreateError(context.TODO(), string(_errors.ErrBadParameter), _errors.Trace(), err.Error(), string(_errors.ErrFromClient))
	}
	return nil
}

// ValidateReq validates the request body of a REST API
func ValidateReq(c echo.Context, req interface{}) error {
	// Bind request body
	if err := c.Bind(req); err != nil {
		return _errors.CreateError(context.TODO(), string(_errors.ErrBadParameter), _errors.Trace(), err.Error(), string(_errors.ErrFromClient))
	}
	// Validate the bound struct
	if err := val.Struct(req); err != nil {
		return _errors.CreateError(context.TODO(), string(_errors.ErrBadParameter), _errors.Trace(), err.Error(), string(_errors.ErrFromClient))
	}
	return nil
}

// ValidateRes validates the response body of a REST API
func ValidateRes(c echo.Context, res interface{}) error {
	// Validate the response struct
	if err := val.Struct(res); err != nil {
		return _errors.CreateError(context.TODO(), string(_errors.ErrBadParameter), _errors.Trace(), err.Error(), string(_errors.ErrFromClient))
	}
	return nil
}
