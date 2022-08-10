package stack

import (
	"fmt"
	"runtime"
	"strings"
)

type Trace []uintptr

func (t *Trace) String() string {
	s := strings.Builder{}
	frames := runtime.CallersFrames(*t)
	for {
		f, more := frames.Next()
		s.WriteString(
			fmt.Sprintf("\nat %s:%d\n   --> %s", f.File, f.Line, f.Function),
		)
		if !more {
			return s.String()
		}
	}
}

func Capture() *Trace {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var t Trace = pcs[0:n]
	return &t
}
