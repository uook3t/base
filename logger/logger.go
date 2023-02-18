package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

func Ctx(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if reqID, ok := ctx.Value(KeyRequestID).(string); ok {
		fields[KeyRequestID] = reqID
	}

	return logrus.WithFields(fields)
}