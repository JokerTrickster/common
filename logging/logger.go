package logging

/*
	로깅 초기화 및 로깅 유틸리티 함수
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// InitLogging initializes loggers for info, warning, and error levels
func InitLogging() error {
	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}
	warningFile, err := os.OpenFile("logs/warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open warning log file: %v", err)
	}
	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	infoLogger = log.New(infoFile, "", log.LstdFlags)
	warningLogger = log.New(warningFile, "", log.LstdFlags)
	errorLogger = log.New(errorFile, "", log.LstdFlags)

	return nil
}

// LogInfo logs an info-level message
func LogInfo(logContent interface{}) {
	infoLogger.Printf("%s", getStringFromInterface(logContent))
}

// LogWarning logs a warning-level message
func LogWarning(logContent interface{}) {
	warningLogger.Printf("%s", getStringFromInterface(logContent))
}

// LogError logs an error-level message
func LogError(logContent interface{}) {
	errorLogger.Printf("%s", getStringFromInterface(logContent))
}

// getStringFromInterface converts an interface to a string
func getStringFromInterface(logContent interface{}) string {
	if reflect.Indirect(reflect.ValueOf(logContent)).Kind() == reflect.Struct {
		raw, err := json.Marshal(logContent)
		if err != nil {
			return fmt.Sprintf("%v", logContent)
		}
		return string(raw)
	}
	return fmt.Sprintf("%v", logContent)
}
