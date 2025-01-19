package middleware

import (
	"encoding/json"
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

			// 요청 종료 후 로깅
			response := c.Response()
			latency := time.Since(startTime)

			if err != nil {
				// 핸들러 실행 중 에러 발생 시 에러 로그 출력

				var resError _error.ResError
				if parseErr := json.Unmarshal([]byte(err.Error()), &resError); parseErr == nil {
					// 에러 로그 출력 (구조화된 에러 정보 포함)
					logger.Error(logging.Log{
						Url:          c.Request().URL.Path,
						Method:       c.Request().Method,
						RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
						Latency:      latency.Milliseconds(),
						HttpCode:     response.Status,
						RequestBody:  requestData.Body,
						RequestPath:  requestData.Path,
						RequestQuery: requestData.Query,
						ErrorInfo: &logging.ErrorInfo{
							Msg:       resError.Msg,     // 에러 메시지
							ErrorType: resError.ErrType, // 에러 타입
							Stack:     resError.Trace,   // 에러 스택
							From:      resError.From,    // 에러 발생 위치
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
