package config

import (
	"github.com/spf13/viper"

	logging "github.com/prospik/places_proxy/pkg/logger"
)

var cfg *viper.Viper

const (
	loggerLevel                   = "logger_level"
	loggerStacktraceLevel         = "logger_st_level"
	loggerSentryLevel             = "logger_sentry_level"
	loggerSentryDSN               = "logger_sentry_dsn"
	loggerSentryStacktraceEnabled = "logger_sentry_stacktrace_enabled"

	httpServerURL          = "http_server_url"
	httpServerName         = "http_server_name"
	httpServerReadTimeout  = "http_server_read_timeout"
	httpServerWriteTimeout = "http_server_write_timeout"

	httpClientName         = "http_client_name"
	httpClientReadTimeout  = "http_client_read_timeout"
	httpClientWriteTimeout = "http_client_write_timeout"

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
	cfg.SetDefault(httpServerURL, ":40001")

	_ = cfg.BindEnv(httpServerName, "HTTP_SERVER_NAME")
	cfg.SetDefault(httpServerName, "places_proxy")

	_ = cfg.BindEnv(httpServerReadTimeout, "HTTP_SERVER_READ_TIMEOUT")
	cfg.SetDefault(httpServerReadTimeout, 30)

	_ = cfg.BindEnv(httpServerWriteTimeout, "HTTP_SERVER_WRITE_TIMEOUT")
	cfg.SetDefault(httpServerWriteTimeout, 30)

	_ = cfg.BindEnv(httpClientName, "HTTP_CLIENT_NAME")
	cfg.SetDefault(httpClientName, "places_proxy")

	_ = cfg.BindEnv(httpClientReadTimeout, "HTTP_CLIENT_READ_TIMEOUT")
	cfg.SetDefault(httpClientReadTimeout, 30)

	_ = cfg.BindEnv(httpClientWriteTimeout, "HTTP_CLIENT_WRITE_TIMEOUT")
	cfg.SetDefault(httpClientWriteTimeout, 30)

	_ = cfg.BindEnv(dbURI, "DB_URI")
	cfg.SetDefault(dbURI, "redis://proxydefaultpass@127.0.0.1:50005/0")
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

// ServerConfig http server configuration
type ServerConfig struct {
	Addr         string
	Name         string
	ReadTimeout  int
	WriteTimeout int
}

// NewServerConfig constructor for ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Addr:         cfg.GetString(httpServerURL),
		Name:         cfg.GetString(httpServerName),
		ReadTimeout:  cfg.GetInt(httpServerReadTimeout),
		WriteTimeout: cfg.GetInt(httpServerWriteTimeout),
	}
}

// ClientConfig http server configuration
type ClientConfig struct {
	Name         string
	ReadTimeout  int
	WriteTimeout int
}

// NewClientConfig constructor for ClientConfig
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Name:         cfg.GetString(httpClientName),
		ReadTimeout:  cfg.GetInt(httpClientReadTimeout),
		WriteTimeout: cfg.GetInt(httpClientWriteTimeout),
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
