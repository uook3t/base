package logger

import (
	"fmt"
	"io"
	"os"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
)

func newFileWriter(fileName string, needStdOut bool) io.Writer {
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

	if !needStdOut {
		return fileWriter
	}
	mw := io.MultiWriter(os.Stdout, fileWriter)
	return mw
}
