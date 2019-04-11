package logging

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestAny(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Tag
	}{
		{
			name: "success work",
			args: struct {
				key   string
				value interface{}
			}{key: "test", value: 123},
			want: &Tag{Key: "test", Value: 123},
		},
		{
			name: "with empty string",
			args: struct {
				key   string
				value interface{}
			}{key: "", value: 123},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Any(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Tag
	}{
		{
			name: "success work",
			args: struct{ err error }{err: errors.New("test")},
			want: &Tag{Key: errorKeyTag, Value: errors.New("test")},
		},
		{
			name: "error not set",
			args: struct{ err error }{err: nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Error(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
