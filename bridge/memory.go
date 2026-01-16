package bridge

import (
	"runtime"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// MemoryModule provides access to Go's memory statistics
type MemoryModule struct{}

// GetStats returns memory statistics to JS
func (m *MemoryModule) GetStats(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		obj := vm.NewObject()
		obj.Set("alloc", ms.Alloc)
		obj.Set("totalAlloc", ms.TotalAlloc)
		obj.Set("sys", ms.Sys)
		obj.Set("numGC", ms.NumGC)

		return obj
	}
}

// RegisterMemory injects the memory object and Ptr constructor into the runtime
func RegisterMemory(vm *goja.Runtime, mf *MemoryFactory, el *eventloop.EventLoop) {
	m := &MemoryModule{}

	// Create the goMemory object
	memObj := vm.NewObject()

	// Add stats method
	memObj.Set("stats", m.GetStats(vm))

	// Add makeShared method
	if mf != nil && el != nil {
		memObj.Set("makeShared", func(call goja.FunctionCall) goja.Value {
			name := call.Argument(0).String()
			size := int(call.Argument(1).ToInteger())

			if name == "" || size <= 0 {
				panic(vm.NewTypeError("makeShared requires a name and positive size"))
			}

			segment := mf.MakeShared(name, size)

			buf := vm.NewArrayBuffer(segment.Data)
			u8 := vm.Get("Uint8Array").ToObject(vm)
			tArray, _ := vm.New(u8, vm.ToValue(buf))

			res := vm.NewObject()
			res.Set("buffer", tArray)
			res.Set("mutex", BindMutex(vm, &segment.Mu, el))

			return res
		})
	}

	vm.Set("goMemory", memObj)

	// Ptr constructor: allows wrapping a value to pass by reference in JS
	vm.Set("Ptr", func(call goja.ConstructorCall) *goja.Object {
		obj := vm.NewObject()
		val := call.Argument(0)
		_ = obj.Set("value", val)
		return obj
	})
}
