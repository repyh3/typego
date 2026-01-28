package core

import (
	"github.com/grafana/sobek"
)

// The returned ArrayBuffer is a copy of the original data, so modifications
// in JavaScript will not affect the Go slice.
//
// For shared memory scenarios where modifications should be visible to both
// Go and JavaScript, use MapSharedBuffer instead.
func ToArrayBuffer(vm *sobek.Runtime, data []byte) sobek.Value {
	return vm.ToValue(vm.NewArrayBuffer(data))
}

// The backing memory is shared between Go and JavaScript, meaning modifications
// from either side are immediately visible to the other.
//
// The buffer is registered as a global variable with the given name, accessible
// as a Uint8Array in JavaScript. This is commonly used for inter-worker
// communication and zero-copy data sharing.
//
// Example:
//
//	data := make([]byte, 1024)
//	core.MapSharedBuffer(vm, "sharedBuffer", data)
//	// In JS: sharedBuffer[0] = 42
//	// In Go: data[0] == 42
func MapSharedBuffer(vm *sobek.Runtime, name string, data []byte) {
	buf := vm.NewArrayBuffer(data)
	view := vm.ToValue(vm.Get("Uint8Array")).ToObject(vm)
	typedArray, _ := vm.New(view, vm.ToValue(buf))
	_ = vm.GlobalObject().Set(name, typedArray)
}
