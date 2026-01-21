// Package fmt provides bindings for Go's fmt package.
package fmt

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
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

type Module struct{}

func (f *Module) Println(call goja.FunctionCall) goja.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return goja.Undefined()
}

func (f *Module) Printf(call goja.FunctionCall) goja.Value {
	format := call.Argument(0).String()
	args := make([]interface{}, len(call.Arguments)-1)
	for i, arg := range call.Arguments[1:] {
		args[i] = arg.Export()
	}
	fmt.Printf(format, args...)
	return goja.Undefined()
}

func Register(vm *goja.Runtime) {
	f := &Module{}

	obj := vm.NewObject()
	_ = obj.Set("Println", f.Println)
	_ = obj.Set("Printf", f.Printf)

	_ = vm.Set("__go_fmt__", obj)
}
