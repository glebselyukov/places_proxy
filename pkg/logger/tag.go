package logging

import (
	"fmt"
)

const (
	errorKeyTag   = "error"
	stacktraceTag = "stacktrace"
)

type Tag struct {
	Key   string      `mapstructure:"key"`
	Value interface{} `mapstructure:"value"`
}

func Any(key string, value interface{}) *Tag {
	if key == "" {
		return nil
	}
	return &Tag{
		Key:   key,
		Value: value,
	}
}

func GoStringer(key string, value fmt.GoStringer) *Tag {
	if key == "" {
		return nil
	}
	return &Tag{
		Key:   key,
		Value: value.GoString(),
	}
}

func Stringer(key string, value fmt.Stringer) *Tag {
	if key == "" {
		return nil
	}
	return &Tag{
		Key:   key,
		Value: value.String(),
	}
}

func Error(err error) *Tag {
	if err == nil {
		return nil
	}
	return &Tag{
		Key:   errorKeyTag,
		Value: err,
	}
}
