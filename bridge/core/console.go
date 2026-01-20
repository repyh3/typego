package core

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

// Console implements the JS console object
type Console struct{}

// Log prints arguments to stdout
func (c *Console) Log(call goja.FunctionCall) goja.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Println(strings.Join(args, " "))
	return goja.Undefined()
}

// Error prints arguments to stderr
func (c *Console) Error(call goja.FunctionCall) goja.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Print("Error: ")
	fmt.Println(strings.Join(args, " "))
	return goja.Undefined()
}

// RegisterConsole injects the console object into the runtime
func RegisterConsole(vm *goja.Runtime) {
	c := &Console{}
	obj := vm.NewObject()
	obj.Set("log", c.Log)
	obj.Set("error", c.Error)
	vm.Set("console", obj)
}
