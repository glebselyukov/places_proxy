package logging

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type sentryCore struct {
	client *raven.Client
	level  zapcore.Level
}

func createSentry(dsn string, level zapcore.Level) (core zapcore.Core, err error) {
	client, err := raven.New(dsn)
	if err != nil {
		return
	}
	if client.Tags == nil {
		client.Tags = make(map[string]string)
	}
	core = &sentryCore{
		client: client,
		level:  level,
	}
	return
}

func (c *sentryCore) Enabled(lvl zapcore.Level) bool {
	return c.level <= lvl
}

func (c sentryCore) With(fields []zapcore.Field) zapcore.Core {
	client := *c.client
	c.client = &client
	objectEncoder := new(ravenTagsObjectEncoder)
	for _, field := range fields {
		field.AddTo(objectEncoder)
		c.client.Tags[objectEncoder.Key] = objectEncoder.Value
	}
	return &c
}

func (c *sentryCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return check.AddCore(entry, c)
	}
	return check
}

type ravenTagsObjectEncoder raven.Tag

func (r *ravenTagsObjectEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	newMap := zapcore.NewMapObjectEncoder()
	err := newMap.AddArray(key, marshaler)
	if err != nil {
		return errors.WithStack(err)
	}
	r.Key = key
	r.Value = fmt.Sprintf("%+v", newMap.Fields)
	return nil
}

func (r *ravenTagsObjectEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	newMap := zapcore.NewMapObjectEncoder()
	err := marshaler.MarshalLogObject(newMap)
	if err != nil {
		return errors.WithStack(err)
	}
	r.Key = key
	r.Value = fmt.Sprintf("%+v", newMap.Fields)
	return nil
}

func (r *ravenTagsObjectEncoder) AddBinary(key string, value []byte) {
	r.Key = key
	r.Value = string(value)
}

func (r *ravenTagsObjectEncoder) AddByteString(key string, value []byte) {
	r.Key = key
	r.Value = string(value)
}

func (r *ravenTagsObjectEncoder) AddBool(key string, value bool) {
	r.Key = key
	switch value {
	case true:
		r.Value = "true"
	default:
		r.Value = "false"
	}
}

func (r *ravenTagsObjectEncoder) AddComplex128(key string, value complex128) {
	r.Key = key
	r.Value = fmt.Sprint(value)
}

func (r *ravenTagsObjectEncoder) AddComplex64(key string, value complex64) {
	r.Key = key
	r.Value = fmt.Sprint(value)
}

func (r *ravenTagsObjectEncoder) AddDuration(key string, value time.Duration) {
	r.Key = key
	r.Value = value.String()
}

func (r *ravenTagsObjectEncoder) AddFloat64(key string, value float64) {
	r.Key = key
	r.Value = strconv.FormatFloat(value, 'f', 6, 64)
}

func (r *ravenTagsObjectEncoder) AddFloat32(key string, value float32) {
	r.Key = key
	r.Value = strconv.FormatFloat(float64(value), 'f', 6, 32)

}

func (r *ravenTagsObjectEncoder) AddInt(key string, value int) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddInt64(key string, value int64) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddInt32(key string, value int32) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddInt16(key string, value int16) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddInt8(key string, value int8) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddString(key, value string) {
	r.Key = key
	r.Value = value
}

func (r *ravenTagsObjectEncoder) AddTime(key string, value time.Time) {
	r.Key = key
	r.Value = value.Format(time.RFC3339)
}

func (r *ravenTagsObjectEncoder) AddUint(key string, value uint) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddUint64(key string, value uint64) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddUint32(key string, value uint32) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddUint16(key string, value uint16) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddUint8(key string, value uint8) {
	r.Key = key
	r.Value = fmt.Sprintf("%d", value)
}

func (r *ravenTagsObjectEncoder) AddUintptr(key string, value uintptr) {
	r.Key = key
	r.Value = fmt.Sprint(value)
}

func (r *ravenTagsObjectEncoder) AddReflected(key string, value interface{}) error {
	r.Key = key
	r.Value = fmt.Sprint(value)
	return nil
}

func (r *ravenTagsObjectEncoder) OpenNamespace(key string) {
	r.Key = key
	r.Value = ""
	return
}

func (c *sentryCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	objectEncoder := new(ravenTagsObjectEncoder)
	tags := make(raven.Tags, 0, len(fields))
	var stack *raven.Stacktrace
	for _, field := range fields {
		if field.Key == stacktraceTag {
			st, ok := field.Interface.(*stackTrace)
			if !ok {
				continue
			}
			stack = st.ravenStackTrace
			continue
		}
		field.AddTo(objectEncoder)
		tags = append(tags, (raven.Tag)(*objectEncoder))
	}
	if stack == nil {
		stack = raven.NewStacktrace(3, 3, nil)
	}

	packet := &raven.Packet{
		Message:   entry.Message,
		Timestamp: raven.Timestamp(entry.Time),
		Platform:  fmt.Sprintf("%s: %s", runtime.GOOS, runtime.Version()),
		Interfaces: []raven.Interface{&raven.Exception{
			Value:      entry.Message,
			Type:       entry.Message,
			Stacktrace: stack,
		}},
		Level:  raven.DEBUG, // TODO: prospik: fix convert level from zap to raven and use convertLevelToSeverity func
		Logger: entry.LoggerName,
		Tags:   tags,
	}

	_, errs := c.client.Capture(packet, c.client.Tags)
	return <-errs
}

func convertLevelToSeverity(level zapcore.Level) raven.Severity {
	switch level {
	case zapcore.FatalLevel:
		return raven.FATAL
	case zapcore.ErrorLevel:
		return raven.ERROR
	case zapcore.WarnLevel:
		return raven.WARNING
	case zapcore.InfoLevel:
		return raven.INFO
	case zapcore.DebugLevel:
		return raven.DEBUG
	default:
		return raven.Severity("UNKNOWN")
	}
}

func (c *sentryCore) Sync() error {
	c.client.Close()
	c.client.Wait()
	return nil
}
