package core

import (
	"fmt"

	"github.com/grafana/sobek"
)

type Console struct{}

func (c *Console) Log(call sobek.FunctionCall) sobek.Value {
	// @optimized: Use []interface{} and fmt.Println to avoid string conversion overhead and allocation.
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Println(args...)
	return sobek.Undefined()
}

func (c *Console) Error(call sobek.FunctionCall) sobek.Value {
	// @optimized: Use []interface{} and fmt.Println to avoid string conversion overhead and allocation.
	args := make([]interface{}, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	fmt.Print("Error: ")
	fmt.Println(args...)
	return sobek.Undefined()
}

func RegisterConsole(vm *sobek.Runtime) {
	c := &Console{}
	obj := vm.NewObject()
	_ = obj.Set("log", c.Log)
	_ = obj.Set("error", c.Error)
	_ = vm.Set("console", obj)
}
