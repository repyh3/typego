package core

import (
	"fmt"
	"strings"

	"github.com/grafana/sobek"
)

type Console struct{}

func (c *Console) Log(call sobek.FunctionCall) sobek.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Println(strings.Join(args, " "))
	return sobek.Undefined()
}

func (c *Console) Error(call sobek.FunctionCall) sobek.Value {
	args := make([]string, len(call.Arguments))
	for i, arg := range call.Arguments {
		args[i] = fmt.Sprint(arg.Export())
	}
	fmt.Print("Error: ")
	fmt.Println(strings.Join(args, " "))
	return sobek.Undefined()
}

func RegisterConsole(vm *sobek.Runtime) {
	c := &Console{}
	obj := vm.NewObject()
	_ = obj.Set("log", c.Log)
	_ = obj.Set("error", c.Error)
	_ = vm.Set("console", obj)
}
