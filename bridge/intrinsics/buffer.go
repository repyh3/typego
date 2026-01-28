package intrinsics

import (
	"github.com/grafana/sobek"
)

// Buffer intrinsics provide high-performance byte array operations.

func (r *Registry) BufferAlloc(call sobek.FunctionCall) sobek.Value {
	size := 0
	if len(call.Arguments) > 0 {
		size = int(call.Arguments[0].ToInteger())
	}

	bytes := make([]byte, size)
	buf := r.vm.NewArrayBuffer(bytes)
	u8 := r.vm.Get("Uint8Array").ToObject(r.vm)
	tArray, _ := r.vm.New(u8, r.vm.ToValue(buf))

	return tArray
}

func (r *Registry) BufferFrom(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return r.vm.ToValue([]byte{})
	}

	arg := call.Arguments[0]

	// If it's a string, we default to UTF-8
	if arg.ExportType().String() == "string" {
		str := arg.String()
		bytes := []byte(str)

		buf := r.vm.NewArrayBuffer(bytes)
		u8 := r.vm.Get("Uint8Array").ToObject(r.vm)
		tArray, _ := r.vm.New(u8, r.vm.ToValue(buf))
		return tArray
	}

	// If it's already a TypedArray or ArrayBuffer, we return a new view
	if _, ok := arg.(*sobek.Object); ok {
		// New Uint8Array from existing buffer/array
		u8 := r.vm.Get("Uint8Array").ToObject(r.vm)
		tArray, _ := r.vm.New(u8, arg)
		return tArray
	}

	return r.vm.ToValue([]byte{})
}
