package env

/*
	시간 변환 유틸리티
*/

import "time"

// TimeToEpochMillis converts time.Time to milliseconds since epoch
func TimeToEpochMillis(t time.Time) int64 {
	return t.UnixNano() / 1_000_000
}

// EpochToTime converts epoch seconds to time.Time
func EpochToTime(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

// EpochToTimeMillis converts epoch milliseconds to time.Time
func EpochToTimeMillis(epoch int64) time.Time {
	return time.Unix(epoch/1_000, (epoch%1_000)*1_000_000)
}
