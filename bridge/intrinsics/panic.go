package intrinsics

import "github.com/grafana/sobek"

// Panic implements the global panic() function for JS.
// Usage: panic("message")
func (r *Registry) Panic(call sobek.FunctionCall) sobek.Value {
	msg := "panic: (nil)"
	if len(call.Arguments) > 0 {
		msg = "panic: " + call.Arguments[0].String()
	}
	// We throw a standard JS Error prefixed with "panic: " to mimic Go's output
	panic(r.vm.ToValue(msg))
}

type panicError struct {
	msg string
}

func newPanicError(msg string) *panicError {
	return &panicError{msg: msg}
}

func (e *panicError) Error() string {
	return e.msg
}
