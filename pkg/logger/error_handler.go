package logging

import (
	"github.com/pkg/errors"
)

type LoggerHandler func(log Logger, err error) (ok bool)

func RegisterHandlers(logger Logger, handlers ...LoggerHandler) (err error) {
	zapLogger, ok := logger.(*zapLogger)
	if !ok {
		return errors.New("current logger not support")
	}
	zapLogger.handlers = handlers
	return
}
