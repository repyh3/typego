// Package worker provides the typego:worker module for multi-threading.
package worker

import (
	"github.com/dop251/goja"
	"github.com/repyh/typego/eventloop"
)

// Handle represents a running worker.
type Handle interface {
	PostMessage(msg goja.Value)
	Terminate()
}

// Spawner is a function that creates a new worker.
type Spawner func(scriptPath string, onMessage func(goja.Value)) (Handle, error)

// Register injects the typego:worker module into the runtime.
func Register(vm *goja.Runtime, el *eventloop.EventLoop, spawner Spawner) {
	obj := vm.NewObject()

	_ = obj.Set("Worker", func(call goja.ConstructorCall) *goja.Object {
		scriptPath := call.Argument(0).String()

		workerObj := vm.NewObject()

		onWorkerMessage := func(msg goja.Value) {
			el.RunOnLoop(func() {
				if onMsg := workerObj.Get("onmessage"); onMsg != nil {
					if fn, ok := goja.AssertFunction(onMsg); ok {
						event := vm.NewObject()
						_ = event.Set("data", msg)
						_, _ = fn(workerObj, event)
					}
				}
			})
		}

		handle, err := spawner(scriptPath, onWorkerMessage)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		_ = workerObj.Set("postMessage", func(call goja.FunctionCall) goja.Value {
			msg := call.Argument(0)
			handle.PostMessage(msg)
			return goja.Undefined()
		})

		_ = workerObj.Set("terminate", func(call goja.FunctionCall) goja.Value {
			handle.Terminate()
			return goja.Undefined()
		})

		return workerObj
	})

	_ = vm.Set("__typego_worker__", obj)
}

// RegisterSelf registers the 'self' object in a worker thread.
func RegisterSelf(vm *goja.Runtime, postToParent func(msg goja.Value)) {
	self := vm.GlobalObject()
	_ = vm.Set("self", self)

	_ = self.Set("postMessage", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0)
		postToParent(msg)
		return goja.Undefined()
	})
}
