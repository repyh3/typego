package intrinsics

import (
	"reflect"

	"github.com/grafana/sobek"
)

// Sizeof implements the global sizeof() function.
// Usage: sizeof(obj) -> returns estimated bytes in memory
func (r *Registry) Sizeof(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) == 0 {
		return sobek.Undefined()
	}

	val := call.Arguments[0]
	size := estimateSize(val)

	// Sobek automatically converts Go int64 to JS number
	return r.vm.ToValue(size)
}

func estimateSize(val sobek.Value) int64 {
	if val == nil {
		return 0
	}

	// 1. Export to Go value to analyze
	export := val.Export()
	if export == nil {
		return 0 // null/undefined
	}

	// 2. Use reflect for Go types
	v := reflect.ValueOf(export)
	switch v.Kind() {
	case reflect.Bool:
		return 1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return 8 // Standardize word size
	case reflect.String:
		return int64(len(v.String()) + 16) // Content + Header overhead
	case reflect.Slice:
		// Header + (Len * ElemSize)
		// Deep inspection is expensive, doing shallow for now + capacity
		elem := v.Type().Elem()
		return int64(v.Cap()) * int64(elem.Size())
	case reflect.Struct:
		return int64(v.Type().Size())
	case reflect.Ptr:
		if v.IsNil() {
			return 8
		}
		return 8 + estimateReflectSize(v.Elem())
	default:
		// Fallback for strict Sobek types (Object, Array, etc) not fully exported
		return 64 // Rough estimate for an object reference
	}
}

func estimateReflectSize(v reflect.Value) int64 {
	if !v.IsValid() {
		return 0
	}
	return int64(v.Type().Size())
}
