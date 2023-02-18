package logger

import "time"

const (
	envDebug = "debug"

	KeyRequestID = "request_id"

	logRotateFileNameFmt = "%s.%%Y%%m%%d"
	logRotateTime        = 24 * time.Hour
	logRotateCount       = 7
	logRotateSize        = 100 * 1024 * 1024
)
