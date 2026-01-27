package intrinsics

import "github.com/grafana/sobek"

// It delegates to the current active scope if any.
func (r *Registry) Recover(call sobek.FunctionCall) sobek.Value {
	if r.currentScope != nil {
		return r.currentScope.RecoverJs(call)
	}
	return sobek.Undefined()
}
