package config

import (
	"github.com/spf13/viper"

	"github.com/prospik/places_proxy/pkg/logger"
)

var cfg *viper.Viper

const (
	loggerLevel                   = "logger_level"
	loggerStacktraceLevel         = "logger_st_level"
	loggerSentryLevel             = "logger_sentry_level"
	loggerSentryDSN               = "logger_sentry_dsn"
	loggerSentryStacktraceEnabled = "logger_sentry_stacktrace_enabled"

	httpServerURL = "http_server_url"

	dbURI = "db_uri"
)

func init() {
	cfg = viper.New()
	cfg.AutomaticEnv()

	_ = cfg.BindEnv(loggerLevel, "LOGGER_LEVEL")
	cfg.SetDefault(loggerLevel, "info")

	_ = cfg.BindEnv(loggerStacktraceLevel, "LOGGER_STACKTRACE_LEVEL")
	cfg.SetDefault(loggerStacktraceLevel, "error")

	_ = cfg.BindEnv(loggerSentryLevel, "LOGGER_SENTRY_LEVEL")
	cfg.SetDefault(loggerSentryLevel, "error")

	_ = cfg.BindEnv(loggerSentryDSN, "LOGGER_SENTRY_DSN")
	cfg.SetDefault(loggerSentryDSN, "http://6a21d946805249c1b58dd3037e842a8d@localhost:9000/1")

	_ = cfg.BindEnv(loggerSentryStacktraceEnabled, "LOGGER_SENTRY_STACKTRACE_ENABLED")
	cfg.SetDefault(loggerSentryStacktraceEnabled, true)

	_ = cfg.BindEnv(httpServerURL, "HTTP_SERVER_URL")
	cfg.SetDefault(httpServerURL, ":40010")

	_ = cfg.BindEnv(dbURI, "DB_URI")
	cfg.SetDefault(dbURI, "postgresql://management:123456@localhost:5432/management?sslmode=disable")
}

// LoggerConfig configuration for logger
type LoggerConfig struct {
	Level                   logging.LevelLogging
	StackTraceLevel         logging.LevelLogging
	SentryLevel             logging.LevelLogging
	SentryDSN               string
	SentryStacktraceEnabled bool
}

// NewLoggerConfig constructor for LoggerConfig
func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:                   logging.LevelFromString(cfg.GetString(loggerLevel)),
		StackTraceLevel:         logging.LevelFromString(cfg.GetString(loggerStacktraceLevel)),
		SentryLevel:             logging.LevelFromString(cfg.GetString(loggerSentryLevel)),
		SentryDSN:               cfg.GetString(loggerSentryDSN),
		SentryStacktraceEnabled: cfg.GetBool(loggerSentryStacktraceEnabled),
	}
}

// ServerConfig grpc server configuration
type ServerConfig struct {
	Addr string
}

// NewServerConfig constructor for GRPCConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr: cfg.GetString(httpServerURL),
	}
}

// DatabaseConfig connection database configuration
type DatabaseConfig struct {
	URI string
}

// NewDatabaseConfig constructor for DatabaseConfig
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URI: cfg.GetString(dbURI),
	}
}
