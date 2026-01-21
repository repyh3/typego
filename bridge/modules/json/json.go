// Package json provides bindings for Go's encoding/json package.
package json

import (
	"encoding/json"

	"github.com/dop251/goja"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
)

func init() {
	core.RegisterModule(&jsonModule{})
}

type jsonModule struct{}

func (m *jsonModule) Name() string {
	return "go:encoding/json"
}

func (m *jsonModule) Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

// Register injects the JSON functions into the runtime
func Register(vm *goja.Runtime) {
	obj := vm.NewObject()
	_ = obj.Set("Marshal", marshal(vm))
	_ = obj.Set("Unmarshal", unmarshal(vm))
	_ = obj.Set("Stringify", marshal(vm)) // Alias for JS familiarity
	_ = obj.Set("Parse", unmarshal(vm))   // Alias for JS familiarity
	_ = vm.Set("__go_json__", obj)
}

// marshal converts a JS value to a JSON string
func marshal(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(vm.NewGoError(ErrMissingArgument))
		}

		value := call.Argument(0).Export()

		indent := ""
		if len(call.Arguments) > 1 {
			indentArg := call.Argument(1)
			if !goja.IsUndefined(indentArg) && !goja.IsNull(indentArg) {
				indent = indentArg.String()
			}
		}

		var result []byte
		var err error

		if indent != "" {
			result, err = json.MarshalIndent(value, "", indent)
		} else {
			result, err = json.Marshal(value)
		}

		if err != nil {
			panic(vm.NewGoError(err))
		}

		return vm.ToValue(string(result))
	}
}

// unmarshal parses a JSON string into a JS value
func unmarshal(vm *goja.Runtime) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(vm.NewGoError(ErrMissingArgument))
		}

		jsonStr := call.Argument(0).String()

		var result interface{}
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			panic(vm.NewGoError(err))
		}

		return vm.ToValue(result)
	}
}
