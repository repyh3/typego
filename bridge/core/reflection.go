package core

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/grafana/sobek"
)

type Binding struct {
	Name   string
	Target interface{}
}

// methodInfo holds metadata about an exported method to avoid repeated reflection lookups.
type methodInfo struct {
	Name  string
	Index int
}

// fieldInfo holds metadata about an exported field to avoid repeated reflection lookups.
type fieldInfo struct {
	Index     int
	Name      string
	Anonymous bool
}

// cache for struct method metadata: reflect.Type -> []methodInfo
var typeMethodCache sync.Map

// cache for struct field metadata: reflect.Type -> []fieldInfo
var typeFieldCache sync.Map

// BindStruct exposes a Go struct to JavaScript with full field and method access.
// Supports nested structs (converted recursively) and callback arguments.
func BindStruct(vm *sobek.Runtime, name string, s interface{}) error {
	obj, err := bindValue(vm, reflect.ValueOf(s), make(map[uintptr]sobek.Value))
	if err != nil {
		return err
	}
	return vm.GlobalObject().Set(name, obj)
}

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
	obj := vm.NewObject()

	if v.CanAddr() {
		visited[v.Addr().Pointer()] = obj
	}

	if err := bindStructFields(vm, obj, v, visited); err != nil {
		return nil, err
	}

	bindMethods(vm, obj, v, visited)

	return obj, nil
}

func bindStructFields(vm *sobek.Runtime, obj *sobek.Object, v reflect.Value, visited map[uintptr]sobek.Value) error {
	t := v.Type()

	// @optimized: Use cached field metadata to avoid repeated reflection overhead (NumField, IsExported).
	var fields []fieldInfo
	cached, loaded := typeFieldCache.Load(t)
	if loaded {
		if cachedFields, ok := cached.([]fieldInfo); ok {
			fields = cachedFields
		}
	}

	if fields == nil {
		numFields := t.NumField()
		fields = make([]fieldInfo, 0, numFields)
		for i := 0; i < numFields; i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			fields = append(fields, fieldInfo{
				Index:     i,
				Name:      field.Name,
				Anonymous: field.Anonymous,
			})
		}
		typeFieldCache.Store(t, fields)
	}

	for _, field := range fields {
		fieldVal := v.Field(field.Index)

		// Support for Flattened Embedding (Anonymous Fields)
		if field.Anonymous {
			// Ensure we are working with a struct for diving
			actual := fieldVal
			for actual.Kind() == reflect.Ptr {
				if actual.IsNil() {
					goto skip
				}
				actual = actual.Elem()
			}
			if actual.Kind() == reflect.Struct {
				if err := bindStructFields(vm, obj, actual, visited); err != nil {
					return err
				}
				continue
			}
		}

	skip:
		jsVal, err := bindValue(vm, fieldVal, visited)
		if err != nil {
			return err
		}
		_ = obj.Set(field.Name, jsVal)
	}
	return nil
}

func bindMethods(vm *sobek.Runtime, obj *sobek.Object, v reflect.Value, visited map[uintptr]sobek.Value) {
	var vPtr reflect.Value
	if v.CanAddr() {
		vPtr = v.Addr()
	} else {
		// Create a copy to get an addressable value if needed for methods
		// This allocation is unavoidable if we want to call pointer methods on value types
		vCopy := reflect.New(v.Type())
		vCopy.Elem().Set(v)
		vPtr = vCopy
	}

	tPtr := vPtr.Type()

	// @optimized: Use cached method metadata to avoid repeated reflection overhead (NumMethod, IsExported).
	var methods []methodInfo
	cached, loaded := typeMethodCache.Load(tPtr)
	if loaded {
		if cachedMethods, ok := cached.([]methodInfo); ok {
			methods = cachedMethods
		}
	}

	if methods == nil {
		numMethods := tPtr.NumMethod()
		methods = make([]methodInfo, 0, numMethods)
		for i := 0; i < numMethods; i++ {
			method := tPtr.Method(i)
			if method.IsExported() {
				methods = append(methods, methodInfo{
					Name:  method.Name,
					Index: i,
				})
			}
		}
		typeMethodCache.Store(tPtr, methods)
	}

	for _, m := range methods {
		methodVal := vPtr.Method(m.Index)
		_ = obj.Set(m.Name, createMethodWrapper(vm, methodVal, m.Name, visited))
	}
}

var argsPool = sync.Pool{
	New: func() interface{} {
		s := make([]reflect.Value, 0, 8)
		return &s
	},
}

func createMethodWrapper(vm *sobek.Runtime, methodVal reflect.Value, methodName string, visited map[uintptr]sobek.Value) func(sobek.FunctionCall) sobek.Value {
	methodType := methodVal.Type()
	numIn := methodType.NumIn()

	return func(call sobek.FunctionCall) sobek.Value {
		// Optimize: Use pool for common small argument lists
		var goArgs []reflect.Value
		var pGoArgs *[]reflect.Value

		if numIn <= 8 {
			pGoArgs = argsPool.Get().(*[]reflect.Value)
			goArgs = (*pGoArgs)[:0:8]
			goArgs = goArgs[:numIn]
			defer func() {
				// Clear references before returning to pool
				for i := range goArgs {
					goArgs[i] = reflect.Value{}
				}
				argsPool.Put(pGoArgs)
			}()
		} else {
			goArgs = make([]reflect.Value, numIn)
		}

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

		// Handle (Value, error) pattern common in Go
		if len(results) == 2 && methodType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if !results[1].IsNil() {
				err := results[1].Interface().(error)
				panic(vm.NewGoError(err))
			}
			val, _ := bindValue(vm, results[0], visited)
			return val
		}

		if len(results) == 1 {
			val, _ := bindValue(vm, results[0], visited)
			return val
		}

		// Multiple results: return as array
		retVals := make([]interface{}, len(results))
		for i, r := range results {
			jsVal, _ := bindValue(vm, r, visited)
			retVals[i] = jsVal
		}
		return vm.NewArray(retVals...)
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
	// @optimized: Pre-allocate slice and use NewArray(vals...) to avoid repeated Set calls.
	l := v.Len()
	vals := make([]interface{}, l)
	for i := 0; i < l; i++ {
		elem, err := bindValue(vm, v.Index(i), visited)
		if err != nil {
			return nil, err
		}
		vals[i] = elem
	}
	return vm.NewArray(vals...), nil
}

func bindMap(vm *sobek.Runtime, v reflect.Value, visited map[uintptr]sobek.Value) (sobek.Value, error) {
	obj := vm.NewObject()
	for _, key := range v.MapKeys() {
		var keyStr string
		// @optimized: Avoid Sprintf if key is already a string.
		if key.Kind() == reflect.String {
			keyStr = key.String()
		} else {
			keyStr = fmt.Sprint(key.Interface())
		}

		val, err := bindValue(vm, v.MapIndex(key), visited)
		if err != nil {
			return nil, err
		}
		_ = obj.Set(keyStr, val)
	}
	return obj, nil
}
