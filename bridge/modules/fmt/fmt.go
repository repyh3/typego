// Package fmt provides bindings for Go's fmt package.
package fmt

import (
	"fmt"

	"github.com/grafana/sobek"
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

func (m *fmtModule) Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

type Module struct{}

func (f *Module) Println(call sobek.FunctionCall) sobek.Value {
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return sobek.Undefined()
}

func (f *Module) Printf(call sobek.FunctionCall) sobek.Value {
	format := call.Argument(0).String()
	args := make([]interface{}, len(call.Arguments)-1)
	for i, arg := range call.Arguments[1:] {
		args[i] = arg.Export()
	}
	fmt.Printf(format, args...)
	return sobek.Undefined()
}

func Register(vm *sobek.Runtime) {
	f := &Module{}

	obj := vm.NewObject()
	_ = obj.Set("Println", f.Println)
	_ = obj.Set("Printf", f.Printf)

	_ = vm.Set("__go_fmt__", obj)
}
