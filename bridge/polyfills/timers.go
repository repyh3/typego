package polyfills

import (
	"time"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// EnableTimers injects setTimeout and setInterval globals
func EnableTimers(vm *goja.Runtime, el *eventloop.EventLoop) {
	vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		fn, _ := goja.AssertFunction(call.Argument(0))
		ms := call.Argument(1).ToInteger()
		el.WGAdd(1)
		go func() {
			time.Sleep(time.Duration(ms) * time.Millisecond)
			el.RunOnLoop(func() {
				_, _ = fn(goja.Undefined())
				el.WGDone()
			})
		}()
		return goja.Undefined()
	})

	vm.Set("setInterval", func(call goja.FunctionCall) goja.Value {
		fn, _ := goja.AssertFunction(call.Argument(0))
		ms := call.Argument(1).ToInteger()

		stop := make(chan struct{})

		go func() {
			ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					el.RunOnLoop(func() {
						_, _ = fn(goja.Undefined())
					})
				case <-stop:
					return
				}
			}
		}()

		// Return a simple ID object for clearInterval
		id := vm.NewObject()
		id.Set("__stop__", stop)
		return id
	})

	vm.Set("clearInterval", func(call goja.FunctionCall) goja.Value {
		obj := call.Argument(0).ToObject(vm)
		if obj != nil {
			if ch := obj.Get("__stop__"); ch != nil {
				if stop, ok := ch.Export().(chan struct{}); ok {
					close(stop)
				}
			}
		}
		return goja.Undefined()
	})
}
