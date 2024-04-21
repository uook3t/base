package logger

import "time"

const (
	envDebug     = "DEBUG"
	envLogSdtOut = "LOGSTDOUT"

	KeyRequestID = "request_id"

	logRotateFileNameFmt = "logs/%s.%%Y%%m%%d.log"
	logRotateTime        = 24 * time.Hour
	logRotateCount       = 7
	logRotateSize        = 100 * 1024 * 1024
)
