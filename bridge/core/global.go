package core

import (
	"github.com/dop251/goja"
)

func RegisterGlobals(vm *goja.Runtime) {
	_ = vm.Set("isGoError", func(call goja.FunctionCall) goja.Value {
		arg := call.Argument(0)
		if goja.IsNull(arg) || goja.IsUndefined(arg) {
			return vm.ToValue(false)
		}

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
