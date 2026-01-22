package intrinsics

import (
	"reflect"

	"github.com/grafana/sobek"
)

// Pointers implements ref() and deref() intrinsics.

func (r *Registry) Ref(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return sobek.Undefined()
	}

	val := call.Arguments[0]
	exported := val.Export()

	// Allocate a new box on the Go heap
	// For now, we'll use a map-based or interface-based box to simulate a pointer
	// In Linker 2.0, this would use raw uintptr + pinning.

	ptr := reflect.New(reflect.TypeOf(exported))
	ptr.Elem().Set(reflect.ValueOf(exported))

	// Return a JS object that mimics the Ref interface
	// { value: T, ptr: uintptr }
	obj := r.vm.NewObject()
	_ = obj.Set("ptr", r.vm.ToValue(ptr.Pointer()))

	// Define 'value' property with getter/setter
	_ = obj.DefineDataProperty("value", r.vm.ToValue(exported), sobek.FLAG_TRUE, sobek.FLAG_TRUE, sobek.FLAG_TRUE)

	return obj
}

func (r *Registry) Deref(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return sobek.Undefined()
	}

	// Simple implementation for now: if it's a Ref object, return its .value
	val := call.Arguments[0]
	if obj := val.ToObject(r.vm); obj != nil {
		if v := obj.Get("value"); v != nil {
			return v
		}
	}

	return sobek.Undefined()
}
