package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log      *zap.Logger
	sugar    *zap.SugaredLogger
	handlers []LoggerHandler
}

func (l zapLogger) Copy(tags ...*Tag) Logger {
	l.log = l.log.With(convertTagsToZapField(tags)...)
	return &l
}

type stdOutEnabler zapcore.Level

func (lvl stdOutEnabler) Enabled(level zapcore.Level) bool {
	return level < zapcore.ErrorLevel && level >= zapcore.Level(lvl)
}

type stdErrEnabler zapcore.Level

func (stdErrEnabler) Enabled(level zapcore.Level) bool {
	return level >= zapcore.ErrorLevel
}

func (l *zapLogger) Write(err error, tags ...*Tag) {
	if err == nil {
		return
	}
	log := l.Copy(tags...)
	var ok bool
	for _, handler := range l.handlers {
		ok = handler(log, err)
		if ok {
			break
		}
	}
	if !ok {
		l.Error(err, tags...)
	}
}

func (l *zapLogger) Fatal(input interface{}, tags ...*Tag) {
	l.send(input, l.log.Fatal, true, tags...)
}

func (l *zapLogger) Error(input interface{}, tags ...*Tag) {
	l.send(input, l.log.Error, true, tags...)
}

func (l *zapLogger) Warning(input interface{}, tags ...*Tag) {
	l.send(input, l.log.Warn, true, tags...)
}

func (l *zapLogger) Info(input interface{}, tags ...*Tag) {
	l.send(input, l.log.Info, false, tags...)
}

// Infow is sugared logging with output only to stdout
func (l *zapLogger) Infow(input string, tags ...interface{}) {
	l.sugar.Infow(input, tags...)
}

func (l *zapLogger) Debug(input interface{}, tags ...*Tag) {
	l.send(input, l.log.Debug, true, tags...)
}

type printingFunc func(msg string, tags ...zap.Field)

func (l *zapLogger) send(input interface{}, printFunc printingFunc, withStack bool, tags ...*Tag) {
	if input == nil {
		l.log.Warn("catch nil input arg in logging")
		return
	}

	var msg string
	switch in := input.(type) {
	case error:
		msg = in.Error()
	case string:
		msg = in
	default:
		msg = fmt.Sprintf("%+v", input)
	}

	if withStack {
		stackTrace := getStackTrace(input)
		tags = append(tags, Any(stacktraceTag, stackTrace))
	}

	printFunc(msg, convertTagsToZapField(tags)...)
}

func convertTagsToZapField(tags []*Tag) []zapcore.Field {
	fields := make([]zap.Field, 0, len(tags))
	for _, tag := range tags {
		fields = append(fields, zap.Any(tag.Key, tag.Value))
	}
	return fields
}
