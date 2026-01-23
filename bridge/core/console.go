package core

import (
	"fmt"
	"sync"

	"github.com/grafana/sobek"
)

var consoleArgsPool = sync.Pool{
	New: func() interface{} {
		s := make([]interface{}, 0, 8)
		return &s
	},
}

type Console struct{}

func (c *Console) Log(call sobek.FunctionCall) sobek.Value {
	// @optimized: Reuse argument slice to reduce allocation frequency.
	args, pArgs := getConsoleArgs(call)
	fmt.Println(args...)
	putConsoleArgs(args, pArgs)
	return sobek.Undefined()
}

func (c *Console) Error(call sobek.FunctionCall) sobek.Value {
	// @optimized: Reuse argument slice to reduce allocation frequency.
	args, pArgs := getConsoleArgs(call)
	fmt.Print("Error: ")
	fmt.Println(args...)
	putConsoleArgs(args, pArgs)
	return sobek.Undefined()
}

func getConsoleArgs(call sobek.FunctionCall) ([]interface{}, *[]interface{}) {
	n := len(call.Arguments)
	var args []interface{}
	var pArgs *[]interface{}

	if n <= 8 {
		pArgs = consoleArgsPool.Get().(*[]interface{})
		args = (*pArgs)[:n]
	} else {
		args = make([]interface{}, n)
	}

	for i, arg := range call.Arguments {
		args[i] = arg.Export()
	}
	return args, pArgs
}

func putConsoleArgs(args []interface{}, pArgs *[]interface{}) {
	if pArgs != nil {
		for i := range args {
			args[i] = nil // Avoid memory leaks
		}
		consoleArgsPool.Put(pArgs)
	}
}

func RegisterConsole(vm *sobek.Runtime) {
	c := &Console{}
	obj := vm.NewObject()
	_ = obj.Set("log", c.Log)
	_ = obj.Set("error", c.Error)
	_ = vm.Set("console", obj)
}
