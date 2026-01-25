package intrinsics

import (
	"reflect"
	"runtime"

	"github.com/grafana/sobek"
)

func (r *Registry) Ref(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return sobek.Undefined()
	}

	val := call.Arguments[0]
	exported := val.Export()

	// Allocate a new box on the Go heap
	ptr := reflect.New(reflect.TypeOf(exported))
	ptr.Elem().Set(reflect.ValueOf(exported))

	obj := r.vm.NewObject()
	_ = obj.Set("ptr", r.vm.ToValue(ptr.Pointer()))

	// Use an accessor to allow direct mutation
	_ = obj.DefineAccessorProperty("value",
		r.vm.ToValue(func(call sobek.FunctionCall) sobek.Value {
			return r.vm.ToValue(ptr.Elem().Interface())
		}),
		r.vm.ToValue(func(call sobek.FunctionCall) sobek.Value {
			newVal := call.Argument(0).Export()
			if newVal == nil {
				ptr.Elem().Set(reflect.Zero(ptr.Elem().Type()))
			} else {
				// Type-safe set if possible
				v := reflect.ValueOf(newVal)
				if v.Type().AssignableTo(ptr.Elem().Type()) {
					ptr.Elem().Set(v)
				}
			}
			return sobek.Undefined()
		}),
		sobek.FLAG_FALSE, sobek.FLAG_TRUE)

	// Finalizer for the handle
	runtime.SetFinalizer(obj, func(o *sobek.Object) {
		// Optional: Log or track cleanup
	})

	return obj
}

func (r *Registry) Deref(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return sobek.Undefined()
	}

	val := call.Arguments[0]
	
	// If it's a number, it's a raw pointer (unsafe)
	if ptr, ok := val.Export().(int64); ok {
		// UNIMPLEMENTED: Raw pointer dereference requires unsafe.Pointer
		// For now, only Ref objects are supported.
		_ = ptr
	}

	// If it's a Ref object
	if obj := val.ToObject(r.vm); obj != nil {
		if v := obj.Get("value"); v != nil {
			return v
		}
	}

	return sobek.Undefined()
}
