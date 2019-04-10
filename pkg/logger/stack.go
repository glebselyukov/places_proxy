package logging

import (
	"fmt"
	"runtime"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type stackTrace struct {
	frames          *runtime.Frames
	ravenStackTrace *raven.Stacktrace
	str             string
}

func (st *stackTrace) String() string {
	if len(st.str) > 0 {
		return st.str
	}
	if st.frames == nil {
		return ""
	}
	st.ravenStackTrace = new(raven.Stacktrace)
	for {
		frame, more := st.frames.Next()
		st.str += fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		st.ravenStackTrace.Frames = append(st.ravenStackTrace.Frames, &raven.StacktraceFrame{
			Function: frame.Function,
			Filename: frame.File,
			Lineno:   frame.Line,
			InApp:    true,
		})
		if !more {
			break
		}
		st.str += "\n"
	}
	return st.str
}

func getStackTrace(input interface{}) *stackTrace {
	var frames []errors.Frame
	err, ok := input.(stackTracer)
	if ok {
		frames = err.StackTrace()
	}
	callers := createPCs(frames)
	return &stackTrace{frames: runtime.CallersFrames(callers)}
}

func createPCs(frames []errors.Frame) []uintptr {
	if len(frames) == 0 {
		var pcs [32]uintptr
		n := runtime.Callers(5, pcs[:])
		return pcs[:n-2]
	}
	callers := make([]uintptr, 0, len(frames))
	for _, st := range frames {
		callers = append(callers, uintptr(st))
	}
	return callers
}
