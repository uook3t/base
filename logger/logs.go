package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Init(serverName string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv(envDebug) == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	var needStdOut bool
	if os.Getenv(envLogSdtOut) == "true" {
		needStdOut = true
	}

	logrus.SetOutput(newFileWriter(serverName, needStdOut))
}
