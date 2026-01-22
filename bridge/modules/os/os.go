// Package os provides bindings for Go's os package.
package os

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
)

func init() {
	core.RegisterModule(&osModule{})
}

type osModule struct{}

func (m *osModule) Name() string {
	return "go:os"
}

func (m *osModule) Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

type Module struct {
	Root string
}

// sanitizePath ensures the path is within the root directory and resolves symlinks
func (m *Module) sanitizePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	realPath, err := filepath.EvalSymlinks(abs)
	if err != nil {
		if os.IsNotExist(err) {
			parentDir := filepath.Dir(abs)
			realParent, pErr := filepath.EvalSymlinks(parentDir)
			if pErr == nil {
				realPath = filepath.Join(realParent, filepath.Base(abs))
			} else {
				return "", pErr
			}
		} else {
			return "", err
		}
	}

	rel, err := filepath.Rel(m.Root, realPath)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return "", os.ErrPermission
	}

	return realPath, nil
}

func (m *Module) WriteFile(vm *sobek.Runtime) func(sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()
		data := call.Argument(1).String()

		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}

		err = os.WriteFile(safePath, []byte(data), 0644)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("os.WriteFile error: %v", err)))
		}

		return sobek.Undefined()
	}
}

func (m *Module) ReadFile(vm *sobek.Runtime) func(sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()

		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}

		data, err := os.ReadFile(safePath)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("os.ReadFile error: %v", err)))
		}

		return vm.ToValue(string(data))
	}
}

func Register(vm *sobek.Runtime) {
	wd, _ := os.Getwd()
	absRoot, _ := filepath.Abs(wd)
	m := &Module{Root: absRoot}

	obj := vm.NewObject()
	_ = obj.Set("WriteFile", m.WriteFile(vm))
	_ = obj.Set("ReadFile", m.ReadFile(vm))

	_ = obj.Set("Getenv", func(call sobek.FunctionCall) sobek.Value {
		key := call.Argument(0).String()
		return vm.ToValue(os.Getenv(key))
	})

	_ = obj.Set("LookupEnv", func(call sobek.FunctionCall) sobek.Value {
		key := call.Argument(0).String()
		val, ok := os.LookupEnv(key)
		result := vm.NewObject()
		_ = result.Set("value", val)
		_ = result.Set("ok", ok)
		return result
	})

	_ = obj.Set("Exit", func(call sobek.FunctionCall) sobek.Value {
		code := int(call.Argument(0).ToInteger())
		os.Exit(code)
		return sobek.Undefined()
	})

	_ = obj.Set("Args", vm.ToValue(os.Args))

	_ = obj.Set("Cwd", func(call sobek.FunctionCall) sobek.Value {
		wd, err := os.Getwd()
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(wd)
	})

	_ = obj.Set("Mkdir", func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()
		if err := os.Mkdir(path, 0755); err != nil {
			panic(vm.NewGoError(err))
		}
		return sobek.Undefined()
	})

	_ = obj.Set("MkdirAll", func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(vm.NewGoError(err))
		}
		return sobek.Undefined()
	})

	_ = obj.Set("Remove", func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()
		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}
		if err := os.Remove(safePath); err != nil {
			panic(vm.NewGoError(err))
		}
		return sobek.Undefined()
	})

	_ = obj.Set("RemoveAll", func(call sobek.FunctionCall) sobek.Value {
		path := call.Argument(0).String()
		safePath, err := m.sanitizePath(path)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("sandbox violation: %v", err)))
		}
		if err := os.RemoveAll(safePath); err != nil {
			panic(vm.NewGoError(err))
		}
		return sobek.Undefined()
	})

	_ = vm.Set("__go_os__", obj)
}
