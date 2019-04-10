package logging

import (
	"reflect"
	"testing"
)

func TestLevelFromString(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want LevelLogging
	}{
		{
			name: "success with fatal level",
			args: struct{ level string }{level: "fatal"},
			want: FatalLevel,
		},
		{
			name: "success with error level",
			args: struct{ level string }{level: "error"},
			want: ErrorLevel,
		},
		{
			name: "success with warning level",
			args: struct{ level string }{level: "warning"},
			want: WarningLevel,
		},
		{
			name: "success with warn level",
			args: struct{ level string }{level: "warn"},
			want: WarningLevel,
		},
		{
			name: "success with info level",
			args: struct{ level string }{level: "info"},
			want: InfoLevel,
		},
		{
			name: "success with debug level",
			args: struct{ level string }{level: "debug"},
			want: DebugLevel,
		},
		{
			name: "success in the different register symbols",
			args: struct{ level string }{level: "FaTal"},
			want: FatalLevel,
		},
		{
			name: "not support level",
			args: struct{ level string }{level: "test"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LevelFromString(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
