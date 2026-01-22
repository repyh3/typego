package core

import (
	"fmt"
	"reflect"

	"github.com/grafana/sobek"
)

// Binding represents a Go struct that has been bound to the JavaScript runtime.
type Binding struct {
	Name   string
	Target interface{}
}

// BindStruct exposes a Go struct to JavaScript with full field and method access.
// Supports nested structs (converted recursively) and callback arguments.
func BindStruct(vm *sobek.Runtime, name string, s interface{}) error {
	obj, err := bindValue(vm, reflect.ValueOf(s), make(map[uintptr]sobek.Value))
	if err != nil {
		return err
	}
	return vm.GlobalObject().Set(name, obj)
}

// bindValue recursively converts a Go value to a JavaScript value.
// The visited map prevents infinite loops for circular references.
func bindValue(vm *sobek.Runtime, v reflect.Value, visited map[uintptr]sobek.Value) (sobek.Value, error) {
	if !v.IsValid() {
		return sobek.Undefined(), nil
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return sobek.Null(), nil
		}
		ptr := v.Pointer()
		if cached, ok := visited[ptr]; ok {
			return cached, nil
		}
		return bindValue(vm, v.Elem(), visited)
	}

	switch v.Kind() {
	case reflect.Struct:
		return bindStruct(vm, v, visited)
	case reflect.Slice, reflect.Array:
		return bindSlice(vm, v, visited)
	case reflect.Map:
		return bindMap(vm, v, visited)
	case reflect.Func:
		return vm.ToValue(v.Interface()), nil
	default:
		return vm.ToValue(v.Interface()), nil
	}
}

func bindStruct(vm *sobek.Runtime, v reflect.Value, visited map[uintptr]sobek.Value) (sobek.Value, error) {
	t := v.Type()
	obj := vm.NewObject()

	if v.CanAddr() {
		visited[v.Addr().Pointer()] = obj
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldVal := v.Field(i)
		jsVal, err := bindValue(vm, fieldVal, visited)
		if err != nil {
			return nil, err
		}
		_ = obj.Set(field.Name, jsVal)
	}

	bindMethods(vm, obj, v, visited)

	return obj, nil
}

func bindMethods(vm *sobek.Runtime, obj *sobek.Object, v reflect.Value, visited map[uintptr]sobek.Value) {
	var vPtr reflect.Value
	if v.CanAddr() {
		vPtr = v.Addr()
	} else {
		vCopy := reflect.New(v.Type())
		vCopy.Elem().Set(v)
		vPtr = vCopy
	}

	tPtr := vPtr.Type()
	for i := 0; i < tPtr.NumMethod(); i++ {
		method := tPtr.Method(i)
		if !method.IsExported() {
			continue
		}

		methodVal := vPtr.Method(i)
		methodName := method.Name

		_ = obj.Set(methodName, createMethodWrapper(vm, methodVal, methodName, visited))
	}
}

func createMethodWrapper(vm *sobek.Runtime, methodVal reflect.Value, methodName string, visited map[uintptr]sobek.Value) func(sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		methodType := methodVal.Type()
		numIn := methodType.NumIn()
		goArgs := make([]reflect.Value, numIn)

		for j := 0; j < numIn; j++ {
			argType := methodType.In(j)

			if j < len(call.Arguments) {
				jsArg := call.Arguments[j]
				goArg, err := convertJSToGo(vm, jsArg, argType)
				if err != nil {
					panic(vm.NewTypeError(fmt.Sprintf("Method %s: Argument %d: %v", methodName, j, err)))
				}
				goArgs[j] = goArg
			} else {
				goArgs[j] = reflect.Zero(argType)
			}
		}

		results := methodVal.Call(goArgs)

		if len(results) == 0 {
			return sobek.Undefined()
		}

		if len(results) == 1 {
			jsVal, _ := bindValue(vm, results[0], visited)
			return jsVal
		}

		arr := vm.NewArray()
		for i, r := range results {
			jsVal, _ := bindValue(vm, r, visited)
			_ = arr.Set(fmt.Sprintf("%d", i), jsVal)
		}
		return arr
	}
}

func convertJSToGo(vm *sobek.Runtime, jsVal sobek.Value, goType reflect.Type) (reflect.Value, error) {
	if goType.Kind() == reflect.Func {
		callable, ok := sobek.AssertFunction(jsVal)
		if !ok {
			return reflect.Value{}, fmt.Errorf("expected function, got %T", jsVal.Export())
		}
		return wrapJSCallback(vm, callable, goType), nil
	}

	exported := jsVal.Export()
	if exported == nil {
		return reflect.Zero(goType), nil
	}

	goVal := reflect.ValueOf(exported)

	if goVal.Type().AssignableTo(goType) {
		return goVal, nil
	}

	if goVal.Type().ConvertibleTo(goType) {
		return goVal.Convert(goType), nil
	}

	return reflect.Value{}, fmt.Errorf("expected %s, got %T", goType, exported)
}

func wrapJSCallback(vm *sobek.Runtime, callable sobek.Callable, goType reflect.Type) reflect.Value {
	return reflect.MakeFunc(goType, func(args []reflect.Value) []reflect.Value {
		jsArgs := make([]sobek.Value, len(args))
		for i, arg := range args {
			jsArgs[i] = vm.ToValue(arg.Interface())
		}

		result, err := callable(sobek.Undefined(), jsArgs...)
		if err != nil {
			// Return zero values on error
			numOut := goType.NumOut()
			results := make([]reflect.Value, numOut)
			for i := 0; i < numOut; i++ {
				results[i] = reflect.Zero(goType.Out(i))
			}
			return results
		}

		// Convert result back to Go
		numOut := goType.NumOut()
		if numOut == 0 {
			return nil
		}

		results := make([]reflect.Value, numOut)
		if numOut == 1 {
			exported := result.Export()
			if exported == nil {
				results[0] = reflect.Zero(goType.Out(0))
			} else {
				results[0] = reflect.ValueOf(exported)
			}
		}
		return results
	})
}

func bindSlice(vm *sobek.Runtime, v reflect.Value, visited map[uintptr]sobek.Value) (sobek.Value, error) {
	arr := vm.NewArray()
	for i := 0; i < v.Len(); i++ {
		elem, err := bindValue(vm, v.Index(i), visited)
		if err != nil {
			return nil, err
		}
		_ = arr.Set(fmt.Sprintf("%d", i), elem)
	}
	return arr, nil
}

func bindMap(vm *sobek.Runtime, v reflect.Value, visited map[uintptr]sobek.Value) (sobek.Value, error) {
	obj := vm.NewObject()
	for _, key := range v.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		val, err := bindValue(vm, v.MapIndex(key), visited)
		if err != nil {
			return nil, err
		}
		_ = obj.Set(keyStr, val)
	}
	return obj, nil
}
