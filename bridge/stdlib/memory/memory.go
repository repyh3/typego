// Package memory provides the typego:memory module for shared memory and stats.
package memory

import (
	"runtime"
	"sync" // standard sync

	"github.com/grafana/sobek"
	modulesync "github.com/repyh/typego/bridge/modules/sync" // TypeGo sync module
	"github.com/repyh/typego/eventloop"
)

// SharedSegment represents a named block of memory shared between Go and JS.
type SharedSegment struct {
	Data []byte
	Mu   sync.RWMutex
}

// Factory manages shared memory segments for an engine instance.
type Factory struct {
	segments map[string]*SharedSegment
	mu       sync.Mutex
}

// NewFactory creates a new memory factory.
func NewFactory() *Factory {
	return &Factory{
		segments: make(map[string]*SharedSegment),
	}
}

// MakeShared gets or creates a shared memory segment.
func (f *Factory) MakeShared(name string, size int) *SharedSegment {
	f.mu.Lock()
	defer f.mu.Unlock()

	if s, ok := f.segments[name]; ok {
		return s
	}

	s := &SharedSegment{Data: make([]byte, size)}
	f.segments[name] = s
	return s
}

// Module implements the typego:memory module.
type Module struct {
	Factory *Factory
}

// GetStats returns memory statistics to JS.
func (m *Module) GetStats(vm *sobek.Runtime) func(sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)

		obj := vm.NewObject()
		_ = obj.Set("alloc", ms.Alloc)
		_ = obj.Set("totalAlloc", ms.TotalAlloc)
		_ = obj.Set("sys", ms.Sys)
		_ = obj.Set("numGC", ms.NumGC)

		return obj
	}
}

// Register injects the typego:memory module into the runtime.
func Register(vm *sobek.Runtime, el *eventloop.EventLoop, f *Factory) {
	if f == nil {
		f = NewFactory()
	}
	m := &Module{Factory: f}

	obj := vm.NewObject()
	_ = obj.Set("stats", m.GetStats(vm))

	_ = obj.Set("makeShared", func(call sobek.FunctionCall) sobek.Value {
		name := call.Argument(0).String()
		size := int(call.Argument(1).ToInteger())

		if name == "" || size <= 0 {
			panic(vm.NewTypeError("makeShared requires a name and positive size"))
		}

		segment := f.MakeShared(name, size)

		buf := vm.NewArrayBuffer(segment.Data)
		u8 := vm.Get("Uint8Array").ToObject(vm)
		tArray, _ := vm.New(u8, vm.ToValue(buf))

		res := vm.NewObject()
		_ = res.Set("buffer", tArray)
		// Reuse the production-ready BindMutex from sync module
		_ = res.Set("mutex", modulesync.BindMutex(vm, &segment.Mu, el))

		return res
	})

	// Ptr factory for referencing values
	_ = obj.Set("ptr", func(call sobek.FunctionCall) sobek.Value {
		val := call.Argument(0)
		return val
	})

	_ = vm.Set("__typego_memory__", obj)
}
