package logger

import (
	"fmt"
	"io"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
)

func newFileHook(fileName string) io.Writer {
	fileWriter, err := rotateLogs.New(
		fmt.Sprintf(logRotateFileNameFmt, fileName),
		rotateLogs.WithLinkName(fileName),
		rotateLogs.WithRotationTime(logRotateTime),
		rotateLogs.WithRotationCount(logRotateCount),
		rotateLogs.WithRotationSize(logRotateSize),
	)

	if err != nil {
		panic(err)
	}
	return fileWriter
}
