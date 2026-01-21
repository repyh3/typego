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

func RegisterModule(m Module) {
	modules = append(modules, m)
}

func InitAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	for _, m := range modules {
		m.Register(vm, el)
	}
}

func GetModules() []Module {
	return modules
}
