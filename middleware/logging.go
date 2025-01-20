package middleware

import (
	"fmt"
	"strings"
	"time"

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
			fmt.Println(requestData)
			if parseErr != nil {
				// 요청 데이터 파싱 실패 시 에러 로그 출력
				logger.Error(logging.Log{
					Url:       c.Request().URL.Path,
					Method:    c.Request().Method,
					HttpCode:  400,
					ErrorInfo: parseErr.Error(),
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
	var httpErrStruct HTTPErrorStruct
	// 에러 메시지 분석
	errMessage := err.Error()
	parts := strings.SplitN(errMessage, ", ", 2)
	if len(parts) == 2 {
		// code=... 추출
		fmt.Sscanf(parts[0], "code=%d", &httpErrStruct.Code)

		// message=... 추출
		httpErrStruct.Message = strings.TrimPrefix(parts[1], "message=")
	}

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
		ErrorInfo:    httpErrStruct.Message,
	})
}
