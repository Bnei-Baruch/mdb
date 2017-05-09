package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/stvp/rollbar"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

// panic if err != nil.
// dumps a stack trace to stderr if err satisfies StackTracer
func Must(err error) {
	if err != nil {
		st, ok := err.(StackTracer)
		if ok {
			fmt.Fprintf(os.Stderr, "%s: %+v\n", st, st.StackTrace())
		}
		panic(err)
	}
}

// Convert errors.StackTrace to rollbar.Stack
func ErrorsToRollbarStack(st StackTracer) rollbar.Stack {
	t := st.StackTrace()
	rs := make(rollbar.Stack, len(t))
	for i, f := range t {
		// Program counter as it's computed internally in errors.Frame
		pc := uintptr(f) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			rs[i] = rollbar.Frame{
				Filename: "unknown",
				Method:   "?",
				Line:     0,
			}
			continue
		}

		// symtab info
		file, line := fn.FileLine(pc)
		name := fn.Name()

		// trim compile time GOPATH from file name
		fileWImportPath := trimGOPATH(name, file)

		// Strip only method name from FQN
		idx := strings.LastIndex(name, "/")
		name = name[idx+1:]
		idx = strings.Index(name, ".")
		name = name[idx+1:]

		rs[i] = rollbar.Frame{
			Filename: fileWImportPath,
			Method:   name,
			Line:     line,
		}
	}

	return rs
}

// Taken AS IS from errors pkg since it's not exported there.
// Check out the source code with good comments on https://github.com/pkg/errors/blob/master/stack.go
func trimGOPATH(name, file string) string {
	const sep = "/"
	goal := strings.Count(name, sep) + 2
	i := len(file)
	for n := 0; n < goal; n++ {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}
	file = file[i+len(sep):]
	return file
}
