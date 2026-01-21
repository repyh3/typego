package core

import (
	"github.com/dop251/goja"
)

// RegisterGlobals injects global helper functions into the runtime
func RegisterGlobals(vm *goja.Runtime) {
	_ = vm.Set("isGoError", func(call goja.FunctionCall) goja.Value {
		arg := call.Argument(0)
		if goja.IsNull(arg) || goja.IsUndefined(arg) {
			return vm.ToValue(false)
		}

		// In Goja, Go errors are often wrapped or passed as objects.
		// For now, we basically check if it looks like an error (has message).
		// A strict check would be verifying the underlying Go type.
		obj := arg.ToObject(vm)
		if obj == nil {
			return vm.ToValue(false)
		}

		msg := obj.Get("message")
		if msg != nil && !goja.IsUndefined(msg) {
			return vm.ToValue(true)
		}

		return vm.ToValue(false)
	})
}
