package middleware

import (
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
				handleAndLogError(c, logger, err, requestData, latency)
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

// handleAndLogError handles and logs errors with structured parsing
func handleAndLogError(c echo.Context, logger *logging.Logger, err error, requestData request.RequestData, latency time.Duration) {
	var resError _error.CustomError

	// 에러가 CustomError 타입인지 확인
	if customErr, ok := err.(_error.CustomError); ok {
		resError = customErr
	}

	// 구조화된 에러 로그 출력
	logger.Error(logging.Log{
		Url:          c.Request().URL.Path,
		Method:       c.Request().Method,
		RequestID:    c.Response().Header().Get(echo.HeaderXRequestID),
		Latency:      latency.Milliseconds(),
		HttpCode:     400,
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

// parseErrorMessage parses the error message string into structured components
func parseErrorMessage(message string) (errorType, errorMsg, stack, from string) {
	// 중괄호 제거
	message = strings.Trim(message, "{}")

	// 공백 기준으로 나누기
	parts := strings.Fields(message)
	if len(parts) >= 4 {
		errorType = parts[0]                // 첫 번째 부분: 에러 타입
		errorMsg = parts[1]                 // 두 번째 부분: 에러 메시지
		stack = parts[2]                    // 세 번째 부분: 스택 정보
		from = strings.Join(parts[3:], " ") // 나머지 부분: 추가 정보
	}

	return errorType, errorMsg, stack, from
}
