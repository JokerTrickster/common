package logging

/*
	로깅 미들웨어
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_errors "github.com/JokerTrickster/common/error"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

// LoggerMiddleware logs incoming HTTP requests and responses
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		requestID := random.String(32)
		c.Set("rID", requestID)
		c.Set("startTime", startTime)

		req := c.Request()
		url := req.URL.Path

		// Read request body
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil && err != io.EOF {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var requestBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			requestBody = nil
		}

		err = next(c)
		resCode := c.Response().Status
		logEntry := &Log{}
		logEntry.MakeLog("", url, req.Method, startTime, resCode, requestID, requestBody, c.QueryParams(), extractPathParams(c))

		if err != nil {
			errInfo := _errors.ParseError(err.Error())
			logEntry.MakeErrorLog(ErrorInfo{
				Stack:     errInfo.Trace,
				ErrorType: errInfo.ErrType,
				Msg:       errInfo.Msg,
				From:      errInfo.From,
			})
			LogError(logEntry)
			return echo.NewHTTPError(errInfo.HttpCode, fmt.Sprintf("%s: %s", errInfo.ErrType, errInfo.Msg))
		}

		LogInfo(logEntry)
		return nil
	}
}

// extractPathParams extracts path parameters from an Echo context
func extractPathParams(c echo.Context) map[string]string {
	params := make(map[string]string)
	for _, name := range c.ParamNames() {
		params[name] = c.Param(name)
	}
	return params
}
