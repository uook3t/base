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

	logrus.SetOutput(newFileWriter(serverName))
}
