package intrinsics

import (
	"unicode/utf8"

	"github.com/grafana/sobek"
)

// Encode implements TextEncoder.prototype.encode
func (r *Registry) Encode(call sobek.FunctionCall) sobek.Value {
	var bytes []byte
	if len(call.Arguments) > 0 {
		str := call.Arguments[0].String()
		bytes = []byte(str)
	} else {
		bytes = []byte{}
	}

	// In Sobek, we can create a Uint8Array from a []byte
	buf := r.vm.NewArrayBuffer(bytes)
	u8 := r.vm.Get("Uint8Array").ToObject(r.vm)
	tArray, _ := r.vm.New(u8, r.vm.ToValue(buf))

	return tArray
}

// Decode implements TextDecoder.prototype.decode
func (r *Registry) Decode(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return r.vm.ToValue("")
	}

	arg := call.Arguments[0]
	exported := arg.Export()

	var bytes []byte
	switch v := exported.(type) {
	case []byte:
		bytes = v
	case sobek.ArrayBuffer:

		bytes = v.Bytes()
	default:
		// Fallback for other typed arrays or objects
		if obj, ok := arg.(*sobek.Object); ok {
			if buffer := obj.Get("buffer"); buffer != nil {
				if ab, ok := buffer.Export().(sobek.ArrayBuffer); ok {
					bytes = ab.Bytes()
				}
			}
		}
	}

	if bytes == nil {
		return r.vm.ToValue("")
	}

	if !utf8.Valid(bytes) {
		// Handle non-UTF8? For now just try to convert.
		return r.vm.ToValue(string(bytes))
	}

	return r.vm.ToValue(string(bytes))
}
