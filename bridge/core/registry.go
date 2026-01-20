package core

import (
	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// Module represents a registerable TypeGo module.
type Module interface {
	// Name returns the module's import path (e.g., "go:fmt", "typego:memory")
	Name() string
	// Register injects the module's bindings into the JavaScript runtime
	Register(vm *goja.Runtime, el *eventloop.EventLoop)
}

var modules []Module

// RegisterModule adds a module to the global registry.
// Modules typically call this in their init() function.
func RegisterModule(m Module) {
	modules = append(modules, m)
}

// InitAll initializes all registered modules.
// Called once during engine startup.
func InitAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	for _, m := range modules {
		m.Register(vm, el)
	}
}

// GetModules returns all registered modules (for debugging/introspection).
func GetModules() []Module {
	return modules
}
