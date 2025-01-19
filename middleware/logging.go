package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_error "github.com/JokerTrickster/common/error"
	"github.com/JokerTrickster/common/logging"
	"github.com/JokerTrickster/common/request"

	"github.com/labstack/echo/v4"
)

// HTTPErrorStruct represents the structure of the error message
type HTTPErrorStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LoggingMiddleware logs the details of each request, including parsed request data and errors
func LoggingMiddleware(logger *logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			// 요청 데이터 파싱
			requestData, parseErr := request.ParseRequest(c)
			if parseErr != nil {
				// 요청 데이터 파싱 실패 시 에러 로그 출력
				logger.Error(logging.Log{
					Url:      c.Request().URL.Path,
					Method:   c.Request().Method,
					HttpCode: 400,
					ErrorInfo: &logging.ErrorInfo{
						Msg:       "Failed to parse request",
						ErrorType: "ParseError",
					},
				})
				return c.JSON(400, map[string]string{"error": "Invalid request format"})
			}

			// 다음 핸들러 실행
			err := next(c)
			fmt.Println(err)

			// 요청 종료 후 로깅
			response := c.Response()
			latency := time.Since(startTime)

			if err != nil {
				// 핸들러 실행 중 에러 발생 시 에러 로그 출력
				var resError _error.ResError
				var httpErrStruct HTTPErrorStruct
				// 에러 메시지를 분리하여 파싱
				errMessage := err.Error()
				parts := strings.SplitN(errMessage, ", ", 2)
				if len(parts) == 2 {
					// code=... 추출
					fmt.Sscanf(parts[0], "code=%d", &httpErrStruct.Code)

					// message=... 추출
					messagePart := strings.TrimPrefix(parts[1], "message=")
					httpErrStruct.Message = messagePart
				}
				fmt.Println("해보자 : ", httpErrStruct)
				// message가 JSON이면 파싱
				if json.Unmarshal([]byte(httpErrStruct.Message), &resError) == nil {
					// 구조화된 에러 로그 출력
					logger.Error(logging.Log{
						Url:          c.Request().URL.Path,
						Method:       c.Request().Method,
						RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
						Latency:      latency.Milliseconds(),
						HttpCode:     httpErrStruct.Code,
						RequestBody:  requestData.Body,
						RequestPath:  requestData.Path,
						RequestQuery: requestData.Query,
						ErrorInfo: &logging.ErrorInfo{
							Msg:       resError.Msg,
							ErrorType: resError.ErrType,
							Stack:     resError.Trace,
							From:      resError.From,
						},
					})
				}
				return err
			}

			// 정상 요청 정보 로깅
			logger.Info(logging.Log{
				Url:          c.Request().URL.Path,
				Method:       c.Request().Method,
				RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
				Latency:      latency.Milliseconds(),
				HttpCode:     response.Status,
				RequestBody:  requestData.Body,
				RequestPath:  requestData.Path,
				RequestQuery: requestData.Query,
			})

			return nil
		}
	}
}
