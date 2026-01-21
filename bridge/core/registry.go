package core

import (
	"github.com/dop251/goja"
	"github.com/repyh/typego/eventloop"
)

type Module interface {
	// Name returns the module's import path (e.g., "go:fmt", "typego:memory")
	Name() string
	Register(vm *goja.Runtime, el *eventloop.EventLoop)
}

var modules []Module

// Modules typically call this in their init() function.
func RegisterModule(m Module) {
	modules = append(modules, m)
}

// Called once during engine startup.
func InitAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	for _, m := range modules {
		m.Register(vm, el)
	}
}

func GetModules() []Module {
	return modules
}
