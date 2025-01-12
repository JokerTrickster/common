package mysql

import (
	"os"
	"time"

	"github.com/google/uuid"
)

// getEnvOrFallback fetches an environment variable or returns a fallback value
func getEnvOrFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// PKIDGenerate generates a UUID for primary keys
func PKIDGenerate() string {
	return uuid.New().String()
}

// NowDateGenerate generates the current date-time string
func NowDateGenerate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// EpochToTime converts epoch seconds to time.Time
func EpochToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

// EpochToTimeString converts epoch seconds to a string
func EpochToTimeString(t int64) string {
	return EpochToTime(t).String()
}

// TimeStringToEpoch converts a time string to epoch seconds
func TimeStringToEpoch(t string) int64 {
	date, _ := time.Parse("2006-01-02 15:04:05", t)
	return date.Unix()
}

// TimeToEpoch converts time.Time to epoch seconds
func TimeToEpoch(t time.Time) int64 {
	return t.Unix()
}
