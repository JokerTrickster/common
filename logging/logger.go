package logging

import (
	"encoding/json"
	"log"
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
	ErrorInfo    string                 `json:"errorInfo,omitempty"`
}

// ErrorInfo contains details about an error
type ErrorInfo struct {
	Stack     string `json:"stack,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Msg       string `json:"msg,omitempty"`
	From      string `json:"from,omitempty"`
}

// Logger is the common logging interface
type Logger struct {
	project string
	env     string
}

// NewLogger initializes a new logger
func NewLogger(project, env string) *Logger {
	log.SetFlags(0)
	return &Logger{
		project: project,
		env:     env,
	}
}

// Info logs informational messages
func (l *Logger) Info(logData Log) {
	logData.Project = l.project
	logData.Env = l.env
	logData.Created = time.Now().Format(time.RFC3339)
	logData.Type = "info"

	l.log(logData)
}

// Error logs error messages
func (l *Logger) Error(logData Log) {
	logData.Project = l.project
	logData.Env = l.env
	logData.Created = time.Now().Format(time.RFC3339)
	logData.Type = "error"

	l.log(logData)
}

// log writes the log entry as JSON to stdout
func (l *Logger) log(logData Log) {
	jsonData, err := json.Marshal(logData)
	if err != nil {
		log.Printf("Failed to marshal log data: %v", err)
		return
	}
	// Print to stdout (CloudWatch Logs will capture this)
	log.Println(string(jsonData))
}
