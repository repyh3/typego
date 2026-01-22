package intrinsics

import (
	"fmt"

	"github.com/grafana/sobek"
)

type scopeState struct {
	defers      []sobek.Callable
	activePanic sobek.Value
	vm          *sobek.Runtime
}

func (s *scopeState) DeferJs(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) > 0 {
		if cb, ok := sobek.AssertFunction(call.Arguments[0]); ok {
			s.defers = append(s.defers, cb)
		}
	}
	return sobek.Undefined()
}

func (s *scopeState) RecoverJs(call sobek.FunctionCall) sobek.Value {
	if s.activePanic != nil {
		p := s.activePanic
		s.activePanic = nil // Reset panic state
		return p
	}
	return sobek.Undefined()
}

// Scope implements the typego.scope() bridge.
// Usage: typego.scope(func(defer, recover) { ... })
func (r *Registry) Scope(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		panic(r.vm.NewGoError(newPanicError("scope requires a callback function")))
	}

	fn, ok := sobek.AssertFunction(call.Arguments[0])
	if !ok {
		panic(r.vm.NewGoError(newPanicError("scope argument 1 must be a function")))
	}

	state := &scopeState{vm: r.vm}

	// Save previous scope for nesting
	prevScope := r.currentScope
	r.currentScope = state

	// Execute defers on Go return or panic
	defer func() {
		if rGo := recover(); rGo != nil {
			// Extract value from Exception
			type valueHolder interface {
				Value() sobek.Value
			}
			if vh, ok := rGo.(valueHolder); ok {
				state.activePanic = vh.Value()
			} else if goErr, ok := rGo.(error); ok {
				state.activePanic = r.vm.ToValue(goErr.Error())
			} else {
				state.activePanic = r.vm.ToValue(fmt.Sprint(rGo))
			}
		}

		// Run JS defers in LIFO order
		for i := len(state.defers) - 1; i >= 0; i-- {
			_, err := state.defers[i](sobek.Undefined())
			if err != nil {
				// A panic in a defer overrides the current panic (standard Go behavior)
				state.activePanic = r.vm.ToValue(err)
			}
		}

		// Re-panic if not recovered
		if state.activePanic != nil {
			r.currentScope = prevScope // Restore BEFORE re-panicking
			panic(state.activePanic)
		} else {
			r.currentScope = prevScope // Restore on normal exit or recovery
		}
	}()

	// Call the user's function, passing our 'defer' and 'recover' functions
	ret, err := fn(sobek.Undefined(), r.vm.ToValue(state.DeferJs), r.vm.ToValue(state.RecoverJs))
	if err != nil {
		// This handles non-panic errors by converting them to panics within the scope
		panic(err)
	}

	return ret
}
