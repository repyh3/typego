package intrinsics

import "github.com/grafana/sobek"

// Recover implements the global recover() function.
// It delegates to the current active scope if any.
func (r *Registry) Recover(call sobek.FunctionCall) sobek.Value {
	if r.currentScope != nil {
		return r.currentScope.RecoverJs(call)
	}
	return sobek.Undefined()
}
