package intrinsics

import (
	"github.com/grafana/sobek"
)

// Scope implements the typego.scope() bridge.
// Usage: typego.scope(func(defer) { ... })
func (r *Registry) Scope(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		panic(r.vm.NewGoError(newPanicError("scope requires a callback function")))
	}

	fn, ok := sobek.AssertFunction(call.Arguments[0])
	if !ok {
		panic(r.vm.NewGoError(newPanicError("scope argument 1 must be a function")))
	}

	// LIFO Stack for defers
	var defers []sobek.Callable

	// The 'defer' function we pass to JS code
	// usage: defer(() => cleanup())
	deferJs := func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) > 0 {
			if cb, ok := sobek.AssertFunction(call.Arguments[0]); ok {
				defers = append(defers, cb)
			}
		}
		return sobek.Undefined()
	}

	// Execute defers on Go return (LIFO)
	defer func() {
		for i := len(defers) - 1; i >= 0; i-- {
			_, _ = defers[i](sobek.Undefined())
		}
	}()

	// Call the user's function, passing our 'defer' function
	ret, err := fn(sobek.Undefined(), r.vm.ToValue(deferJs))
	if err != nil {
		// If JS throws, panic so Sobek propagates it up.
		// Go's defer will still run!
		panic(err)
	}

	return ret
}
