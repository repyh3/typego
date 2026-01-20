// Package os provides bindings for Go's os package.
package os

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge/core"
	"github.com/repyh3/typego/eventloop"
)

func init() {
	core.RegisterModule(&osModule{})
}

type osModule struct{}

func (m *osModule) Name() string {
	return "go:os"
}

func (m *osModule) Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

// Module implements the go:os package bindings.
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

// WriteFile maps to os.WriteFile
func (m *Module) WriteFile(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
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

		return goja.Undefined()
	}
}

// ReadFile maps to os.ReadFile
func (m *Module) ReadFile(vm *goja.Runtime) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
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

// Register injects the os functions into the runtime
func Register(vm *goja.Runtime) {
	wd, _ := os.Getwd()
	absRoot, _ := filepath.Abs(wd)
	m := &Module{Root: absRoot}

	obj := vm.NewObject()
	obj.Set("WriteFile", m.WriteFile(vm))
	obj.Set("ReadFile", m.ReadFile(vm))

	vm.Set("__go_os__", obj)
}
