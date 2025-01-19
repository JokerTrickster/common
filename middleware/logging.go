package middleware

import (
	"log"
	"time"

	"github.com/JokerTrickster/common/logging"
	"github.com/JokerTrickster/common/request"

	"github.com/labstack/echo/v4"
)

// LoggingMiddleware logs the details of each request, including parsed request data
func LoggingMiddleware(logger *logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			// 요청 데이터 파싱
			requestData, err := request.ParseRequest(c)
			if err != nil {
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

			// 로깅: 요청 시작
			request := c.Request()
			log.Printf("Request started: Method=%s, Path=%s", request.Method, request.URL.Path)

			// 다음 핸들러 실행
			err = next(c)

			// 로깅: 요청 종료
			response := c.Response()
			latency := time.Since(startTime)
			logger.Info(logging.Log{
				Url:          request.URL.Path,
				Method:       request.Method,
				Latency:      latency.Milliseconds(),
				HttpCode:     response.Status,
				RequestBody:  requestData.Body,
				RequestPath:  requestData.Path,
				RequestQuery: requestData.Query,
			})

			return err
		}
	}
}
