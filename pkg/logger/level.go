package logging

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type LevelLogging zapcore.Level

const (
	FatalLevel = LevelLogging(iota)
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
)

func LevelFromString(level string) LevelLogging {
	level = strings.ToLower(level)
	switch level {
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warning", "warn":
		return WarningLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return -1
	}
}

func convertToZapLevel(lvl LevelLogging) zapcore.Level {
	switch lvl {
	case FatalLevel:
		return zapcore.FatalLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case WarningLevel:
		return zapcore.WarnLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case DebugLevel:
		return zapcore.DebugLevel
	default:
		return -1
	}
}
