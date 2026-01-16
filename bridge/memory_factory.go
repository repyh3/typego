package bridge

import (
	"sync"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

type SharedSegment struct {
	Data []byte
	Mu   sync.RWMutex
}

type MemoryFactory struct {
	segments map[string]*SharedSegment
	mu       sync.Mutex
}

func NewMemoryFactory() *MemoryFactory {
	return &MemoryFactory{
		segments: make(map[string]*SharedSegment),
	}
}

func (f *MemoryFactory) MakeShared(name string, size int) *SharedSegment {
	f.mu.Lock()
	defer f.mu.Unlock()

	if s, ok := f.segments[name]; ok {
		return s
	}

	s := &SharedSegment{Data: make([]byte, size)}
	f.segments[name] = s
	return s
}

func RegisterMemoryFactory(vm *goja.Runtime, factory *MemoryFactory, el *eventloop.EventLoop) {
	obj := vm.NewObject()

	obj.Set("makeShared", func(call goja.FunctionCall) goja.Value {
		name := call.Argument(0).String()
		size := int(call.Argument(1).ToInteger())

		if name == "" || size <= 0 {
			panic(vm.NewTypeError("makeShared requires a name and positive size"))
		}

		segment := factory.MakeShared(name, size)

		buf := vm.NewArrayBuffer(segment.Data)
		u8 := vm.Get("Uint8Array").ToObject(vm)
		tArray, _ := vm.New(u8, vm.ToValue(buf))

		res := vm.NewObject()
		res.Set("buffer", tArray)
		res.Set("mutex", BindMutex(vm, &segment.Mu, el))

		return res
	})

	_ = vm.Set("__go_memory_factory__", obj)
}

// EnableMemoryFactory exposes makeShared directly on globalThis for standalone binaries
func EnableMemoryFactory(vm *goja.Runtime, sharedBuffers map[string][]byte) {
	vm.Set("makeShared", func(call goja.FunctionCall) goja.Value {
		name := call.Argument(0).String()
		size := int(call.Argument(1).ToInteger())

		if name == "" || size <= 0 {
			panic(vm.NewTypeError("makeShared requires a name and positive size"))
		}

		// Check if already exists
		if buf, ok := sharedBuffers[name]; ok {
			return vm.ToValue(vm.NewArrayBuffer(buf))
		}

		// Create new buffer
		buf := make([]byte, size)
		sharedBuffers[name] = buf
		return vm.ToValue(vm.NewArrayBuffer(buf))
	})
}
