package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Write(err error, tags ...*Tag)

	Fatal(input interface{}, tags ...*Tag)
	Error(input interface{}, tags ...*Tag)
	Warning(input interface{}, tags ...*Tag)
	Info(input interface{}, tags ...*Tag)
	Debug(input interface{}, tags ...*Tag)

	Copy(tags ...*Tag) Logger

	Infow(input string, tags ...interface{})
}

func New(opts ...Opt) (logger Logger) {
	options := newOptions(opts...)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	tee := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdOutEnabler(convertToZapLevel(options.level))),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), stdErrEnabler(0)),
	)

	defer func() {
		log := zap.New(tee)
		sugar := log.Sugar()
		_, err := zap.RedirectStdLogAt(log, zapcore.DebugLevel)
		if err != nil {
			log.Fatal(err.Error())
		}
		logger = &zapLogger{
			log:   log,
			sugar: sugar,
		}
	}()

	if options.sentry.dsn == "" {
		return
	}

	sentryCore, err := createSentry(options.sentry.dsn, convertToZapLevel(options.sentry.level))
	if err != nil {
	}
	tee = zapcore.NewTee(tee, sentryCore)

	return
}
