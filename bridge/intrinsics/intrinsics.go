package intrinsics

import (
	"sync"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/eventloop"
)

type Registry struct {
	vm           *sobek.Runtime
	currentScope *scopeState
	VMLock       sync.Mutex
	el           *eventloop.EventLoop
}

// Enable registers all global intrinsics (panic, sizeof, defer/scope)
func Enable(vm *sobek.Runtime, el *eventloop.EventLoop) *Registry {
	r := &Registry{vm: vm, el: el}

	_ = vm.Set("panic", r.Panic)
	_ = vm.Set("sizeof", r.Sizeof)
	_ = vm.Set("ref", r.Ref)
	_ = vm.Set("deref", r.Deref)
	_ = vm.Set("recover", r.Recover)
	_ = vm.Set("go", r.Go)
	_ = vm.Set("makeChan", r.MakeChan)
	_ = vm.Set("select", r.Select)
	_ = vm.Set("cap", r.Cap)
	_ = vm.Set("make", r.Make)
	_ = vm.Set("copy", r.Copy)
	_ = vm.Set("wrapReader", r.WrapReader)
	_ = vm.Set("wrapWriter", r.WrapWriter)

	// Background backup for polyfills
	_ = vm.Set("__encode", r.Encode)
	_ = vm.Set("__decode", r.Decode)
	_ = vm.Set("__bufferAlloc", r.BufferAlloc)
	_ = vm.Set("__bufferFrom", r.BufferFrom)

	// Scope needs the VM to create the defer callback

	_ = vm.Set("typego", map[string]interface{}{
		"scope": r.Scope,
	})

	r.EnableGlobals()

	return r
}
