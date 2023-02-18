package errors

import (
	stderrors "errors"
	"fmt"
	"github.com/pkg/errors"
)

func New(msg string) error {
	return stderrors.New(msg)
}

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

func Wrap(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args)
}
