// Package fmt provides bindings for Go's fmt package.
package fmt

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge/core"
	"github.com/repyh3/typego/eventloop"
)

func init() {
	core.RegisterModule(&fmtModule{})
}

type fmtModule struct{}

func (m *fmtModule) Name() string {
	return "go:fmt"
}

func (m *fmtModule) Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

// Module implements the go:fmt package bindings.
type Module struct{}

// Println maps to fmt.Println
func (f *Module) Println(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

// Printf maps to fmt.Printf
func (f *Module) Printf(call goja.FunctionCall) goja.Value {
	format := call.Argument(0).String()
	args := make([]interface{}, len(call.Arguments)-1)
	for i, arg := range call.Arguments[1:] {
		args[i] = arg.Export()
	}
	fmt.Printf(format, args...)
	return goja.Undefined()
}

// Register injects the fmt functions into the runtime
func Register(vm *goja.Runtime) {
	f := &Module{}

	obj := vm.NewObject()
	obj.Set("Println", f.Println)
	obj.Set("Printf", f.Printf)

	vm.Set("__go_fmt__", obj)
}
