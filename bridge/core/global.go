package core

import (
	"github.com/grafana/sobek"
)

func RegisterGlobals(vm *sobek.Runtime) {
	_ = vm.Set("isGoError", func(call sobek.FunctionCall) sobek.Value {
		arg := call.Argument(0)
		if sobek.IsNull(arg) || sobek.IsUndefined(arg) {
			return vm.ToValue(false)
		}

		obj := arg.ToObject(vm)
		if obj == nil {
			return vm.ToValue(false)
		}

		msg := obj.Get("message")
		if msg != nil && !sobek.IsUndefined(msg) {
			return vm.ToValue(true)
		}

		return vm.ToValue(false)
	})
}
