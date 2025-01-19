package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	_error "github.com/JokerTrickster/common/error"
	"github.com/JokerTrickster/common/logging"
	"github.com/JokerTrickster/common/request"

	"github.com/labstack/echo/v4"
)

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
			fmt.Println(requestData.Body)
			fmt.Println(requestData.Path)
			fmt.Println(requestData.Query)

			// 요청 종료 후 로깅
			response := c.Response()
			latency := time.Since(startTime)
			if err != nil {
				// Echo HTTPError 처리
				// Echo HTTPError 처리
				if httpErr, ok := err.(*echo.HTTPError); ok {
					var resError _error.ResError

					// Message가 JSON인지 확인하고 파싱
					if messageJSON, ok := httpErr.Message.(string); ok {
						if parseErr := json.Unmarshal([]byte(messageJSON), &resError); parseErr == nil {
							// 에러 로그 출력
							logger.Error(logging.Log{
								Url:          c.Request().URL.Path,
								Method:       c.Request().Method,
								RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
								Latency:      latency.Milliseconds(),
								HttpCode:     httpErr.Code,
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
							return err
						}
					}
				}
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
