package logging

/*
	로깅 포맷 처리 함수
*/

import (
	"time"
)

// Log represents the structure of a log entry
type Log struct {
	Project      string                 `json:"project"`
	Created      string                 `json:"created"`
	Env          string                 `json:"env"`
	Type         string                 `json:"type"`
	UserID       string                 `json:"userID,omitempty"`
	Url          string                 `json:"url"`
	Method       string                 `json:"method"`
	Latency      int64                  `json:"latency"`
	HttpCode     int                    `json:"httpCode"`
	RequestID    string                 `json:"requestID"`
	RequestBody  map[string]interface{} `json:"requestBody,omitempty"`
	RequestPath  map[string]string      `json:"requestPath,omitempty"`
	RequestQuery map[string][]string    `json:"requestQuery,omitempty"`
	ErrorInfo    *ErrorInfo             `json:"errorInfo,omitempty"`
}

// ErrorInfo contains details about an error
type ErrorInfo struct {
	Stack     string `json:"stack,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Msg       string `json:"msg,omitempty"`
	From      string `json:"from,omitempty"`
}

// MakeLog populates a log structure
func (l *Log) MakeLog(userID, url, method string, startTime time.Time, httpCode int, requestID string, requestBody map[string]interface{}, queryParams map[string][]string, pathValues map[string]string) {
	l.Project = "food-recommendation"
	l.Type = "info"
	l.Env = "production" // Adjust based on your environment
	l.UserID = userID
	l.Created = startTime.Format("2006-01-02 15:04:05")
	l.Url = url
	l.Method = method
	l.Latency = time.Since(startTime).Milliseconds()
	l.HttpCode = httpCode
	l.RequestID = requestID
	l.RequestBody = requestBody
	l.RequestQuery = queryParams
	l.RequestPath = pathValues
}

// MakeErrorLog populates error-related details in the log structure
func (l *Log) MakeErrorLog(errInfo ErrorInfo) {
	l.Type = "error"
	l.ErrorInfo = &errInfo
}
