// Package sync provides bindings for Go's sync package and concurrency primitives.
package sync

import (
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
)

func init() {
	core.RegisterModule(&syncModule{})
}

type syncModule struct{}

func (m *syncModule) Name() string {
	return "go:sync"
}

func (m *syncModule) Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	Register(vm, el)
}

type Module struct {
	el *eventloop.EventLoop
}

func (m *Module) Spawn(call goja.FunctionCall) goja.Value {
	fn, ok := goja.AssertFunction(call.Argument(0))
	if !ok {
		panic(m.el.VM.NewTypeError("Spawn expects a function"))
	}

	go func() {
		m.el.RunOnLoop(func() {
			val, err := fn(goja.Undefined())
			if err != nil {
				return
			}

			if obj := val.ToObject(m.el.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !goja.IsUndefined(then) {
					if thenFn, ok := goja.AssertFunction(then); ok {
						m.el.WGAdd(1)
						done := m.el.VM.ToValue(func(goja.FunctionCall) goja.Value {
							m.el.WGDone()
							return goja.Undefined()
						})
						_, _ = thenFn(val, done, done)
					}
				}
			}
		})
	}()

	return goja.Undefined()
}

func (m *Module) Sleep(call goja.FunctionCall) goja.Value {
	ms := call.Argument(0).ToInteger()
	p, resolve, _ := m.el.CreatePromise()

	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		m.el.RunOnLoop(func() {
			resolve(goja.Undefined())
		})
	}()

	return p
}

func Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	m := &Module{el: el}
	obj := vm.NewObject()
	_ = obj.Set("Spawn", m.Spawn)
	_ = obj.Set("Sleep", m.Sleep)

	_ = vm.Set("Chan", func(call goja.ConstructorCall) *goja.Object {
		ch := make(chan goja.Value, 100)
		res := vm.NewObject()
		_ = res.Set("send", func(c goja.FunctionCall) goja.Value {
			ch <- c.Argument(0)
			return goja.Undefined()
		})
		_ = res.Set("recv", func(c goja.FunctionCall) goja.Value {
			p, resolve, _ := el.CreatePromise()
			go func() { resolve(<-ch) }()
			return p
		})
		return res
	})

	_ = vm.Set("__go_sync__", obj)
}

// AsyncMutex wraps a sync.RWMutex with async-friendly locking operations.
type AsyncMutex struct {
	mu *sync.RWMutex
	el *eventloop.EventLoop
}

func NewAsyncMutex(mu *sync.RWMutex, el *eventloop.EventLoop) *AsyncMutex {
	return &AsyncMutex{mu: mu, el: el}
}

func (m *AsyncMutex) Lock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.Lock()
		resolve(goja.Undefined())
	}()
	return p
}

func (m *AsyncMutex) Unlock() {
	m.mu.Unlock()
}

func (m *AsyncMutex) RLock(vm *goja.Runtime) goja.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.RLock()
		resolve(goja.Undefined())
	}()
	return p
}

func (m *AsyncMutex) RUnlock() {
	m.mu.RUnlock()
}

func BindMutex(vm *goja.Runtime, mu *sync.RWMutex, el *eventloop.EventLoop) goja.Value {
	am := NewAsyncMutex(mu, el)
	obj := vm.NewObject()
	_ = obj.Set("lock", am.Lock)
	_ = obj.Set("unlock", am.Unlock)
	_ = obj.Set("rlock", am.RLock)
	_ = obj.Set("runlock", am.RUnlock)
	return obj
}
