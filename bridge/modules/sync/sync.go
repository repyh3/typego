package sync

import (
	"sync"
	"time"

	"github.com/grafana/sobek"
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

func (m *syncModule) Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	Register(vm, el)
}

type Module struct {
	el *eventloop.EventLoop
}

func (m *Module) Spawn(call sobek.FunctionCall) sobek.Value {
	fn, ok := sobek.AssertFunction(call.Argument(0))
	if !ok {
		panic(m.el.VM.NewTypeError("Spawn expects a function"))
	}

	go func() {
		m.el.RunOnLoop(func() {
			val, err := fn(sobek.Undefined())
			if err != nil {
				return
			}

			if obj := val.ToObject(m.el.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !sobek.IsUndefined(then) {
					if thenFn, ok := sobek.AssertFunction(then); ok {
						m.el.WGAdd(1)
						done := m.el.VM.ToValue(func(sobek.FunctionCall) sobek.Value {
							m.el.WGDone()
							return sobek.Undefined()
						})
						_, _ = thenFn(val, done, done)
					}
				}
			}
		})
	}()

	return sobek.Undefined()
}

func (m *Module) Sleep(call sobek.FunctionCall) sobek.Value {
	ms := call.Argument(0).ToInteger()
	p, resolve, _ := m.el.CreatePromise()

	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		m.el.RunOnLoop(func() {
			resolve(sobek.Undefined())
		})
	}()

	return p
}

func Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	m := &Module{el: el}
	obj := vm.NewObject()
	_ = obj.Set("Spawn", m.Spawn)
	_ = obj.Set("Sleep", m.Sleep)

	_ = vm.Set("Chan", func(call sobek.ConstructorCall) *sobek.Object {
		ch := make(chan sobek.Value, 100)
		res := vm.NewObject()
		_ = res.Set("send", func(c sobek.FunctionCall) sobek.Value {
			ch <- c.Argument(0)
			return sobek.Undefined()
		})
		_ = res.Set("recv", func(c sobek.FunctionCall) sobek.Value {
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

func (m *AsyncMutex) Lock(vm *sobek.Runtime) sobek.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.Lock()
		resolve(sobek.Undefined())
	}()
	return p
}

func (m *AsyncMutex) Unlock() {
	m.mu.Unlock()
}

func (m *AsyncMutex) RLock(vm *sobek.Runtime) sobek.Value {
	p, resolve, _ := m.el.CreatePromise()
	go func() {
		m.mu.RLock()
		resolve(sobek.Undefined())
	}()
	return p
}

func (m *AsyncMutex) RUnlock() {
	m.mu.RUnlock()
}

func BindMutex(vm *sobek.Runtime, mu *sync.RWMutex, el *eventloop.EventLoop) sobek.Value {
	am := NewAsyncMutex(mu, el)
	obj := vm.NewObject()
	_ = obj.Set("lock", am.Lock)
	_ = obj.Set("unlock", am.Unlock)
	_ = obj.Set("rlock", am.RLock)
	_ = obj.Set("runlock", am.RUnlock)
	return obj
}
