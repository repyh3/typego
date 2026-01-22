package intrinsics

import (
	"github.com/grafana/sobek"
)

// Cap implements cap(v)
func (r *Registry) Cap(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		return r.vm.ToValue(0)
	}

	val := call.Arguments[0]
	if obj, ok := val.(*sobek.Object); ok {
		// For TypedArrays (Uint8Array, etc.)
		if buffer := obj.Get("buffer"); buffer != nil {
			if bLen := buffer.ToObject(r.vm).Get("byteLength"); bLen != nil {
				if bpe := obj.Get("BYTES_PER_ELEMENT"); bpe != nil {
					// capacity = total bytes / bytes per element
					cap := bLen.ToInteger() / bpe.ToInteger()
					return r.vm.ToValue(cap)
				}
			}
		}

		// Fallback to length for standard Arrays
		if l := obj.Get("length"); l != nil {
			return l
		}
	}

	return r.vm.ToValue(0)
}

// Make implements make(Type, len, cap)
func (r *Registry) Make(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 2 {
		panic(r.vm.NewTypeError("make requires at least a type and a length"))
	}

	typ := call.Arguments[0]
	length := int(call.Arguments[1].ToInteger())
	capacity := length
	if len(call.Arguments) > 2 {
		capacity = int(call.Arguments[2].ToInteger())
	}

	if capacity < length {
		panic(r.vm.NewTypeError("capacity cannot be less than length"))
	}

	// If typ is a constructor (e.g. Uint8Array)
	if _, ok := sobek.AssertFunction(typ); ok {
		// New TypedArray(capacity)
		arr, err := r.vm.New(typ, r.vm.ToValue(capacity))
		if err != nil {
			panic(err)
		}

		// If capacity > length, we need a subarray
		if capacity > length {
			sub, _ := sobek.AssertFunction(arr.Get("subarray"))
			res, _ := sub(arr, r.vm.ToValue(0), r.vm.ToValue(length))
			return res
		}
		return arr
	}

	return sobek.Undefined()
}

// Copy implements copy(dst, src)
func (r *Registry) Copy(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 2 {
		return r.vm.ToValue(0)
	}

	dst := call.Arguments[0].ToObject(r.vm)
	src := call.Arguments[1].ToObject(r.vm)

	set, ok := sobek.AssertFunction(dst.Get("set"))
	if !ok {
		return r.vm.ToValue(0)
	}

	dstLen := int(dst.Get("length").ToInteger())
	srcLen := int(src.Get("length").ToInteger())

	n := dstLen
	if srcLen < n {
		n = srcLen
	}

	if n == 0 {
		return r.vm.ToValue(0)
	}

	var finalSrc sobek.Value = src
	if srcLen > n {
		sub, _ := sobek.AssertFunction(src.Get("subarray"))
		finalSrc, _ = sub(src, r.vm.ToValue(0), r.vm.ToValue(n))
	}

	_, err := set(dst, finalSrc)
	if err != nil {
		return r.vm.ToValue(0)
	}

	return r.vm.ToValue(n)
}
