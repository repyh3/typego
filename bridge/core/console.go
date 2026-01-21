package core

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

type Console struct{}

func (c *Console) Log(call goja.FunctionCall) goja.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Println(strings.Join(args, " "))
	return goja.Undefined()
}

func (c *Console) Error(call goja.FunctionCall) goja.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Print("Error: ")
	fmt.Println(strings.Join(args, " "))
	return goja.Undefined()
}

func RegisterConsole(vm *goja.Runtime) {
	c := &Console{}
	obj := vm.NewObject()
	_ = obj.Set("log", c.Log)
	_ = obj.Set("error", c.Error)
	_ = vm.Set("console", obj)
}
