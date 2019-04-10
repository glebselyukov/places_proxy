package logging

type options struct {
	level           LevelLogging
	stackTraceLevel LevelLogging
	sentry          *sentryOptions
}

type sentryOptions struct {
	dsn                 string
	level               LevelLogging
	isStackTraceEnabled bool
}

func newOptions(opts ...Opt) *options {
	options := &options{
		level:           DebugLevel,
		stackTraceLevel: ErrorLevel,
		sentry: &sentryOptions{
			level:               ErrorLevel,
			isStackTraceEnabled: true,
		},
	}
	for index := range opts {
		opts[index](options)
	}
	return options
}

type Opt func(options *options)

func Level(lvl LevelLogging) Opt {
	return func(options *options) {
		options.level = lvl
	}
}

func SentryDSN(dsn string) Opt {
	return func(options *options) {
		options.sentry.dsn = dsn
	}
}

func StackTraceLevel(lvl LevelLogging) Opt {
	return func(options *options) {
		options.stackTraceLevel = lvl
	}
}

func SentryLevel(lvl LevelLogging) Opt {
	return func(options *options) {
		options.sentry.level = lvl
	}
}

func SentryStacktraceEnabled(isEnabled bool) Opt {
	return func(options *options) {
		options.sentry.isStackTraceEnabled = isEnabled
	}
}
