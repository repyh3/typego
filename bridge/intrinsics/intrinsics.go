package intrinsics

import (
	"github.com/grafana/sobek"
)

// Registry holds references needed for intrinsics
type Registry struct {
	vm *sobek.Runtime
}

// Enable registers all global intrinsics (panic, sizeof, defer/scope)
func Enable(vm *sobek.Runtime) {
	r := &Registry{vm: vm}

	_ = vm.Set("panic", r.Panic)
	_ = vm.Set("sizeof", r.Sizeof)

	// Scope needs the VM to create the defer callback
	_ = vm.Set("typego", map[string]interface{}{
		"scope": r.Scope,
	})
}
