package bridge

import (
	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

type WorkerHandle interface {
	PostMessage(msg goja.Value)
	Terminate()
}

type WorkerSpawner func(scriptPath string, onMessage func(goja.Value)) (WorkerHandle, error)

func RegisterWorker(vm *goja.Runtime, el *eventloop.EventLoop, spawner WorkerSpawner) {
	vm.Set("Worker", func(call goja.ConstructorCall) *goja.Object {
		scriptPath := call.Argument(0).String()

		obj := vm.NewObject()

		onWorkerMessage := func(msg goja.Value) {
			el.RunOnLoop(func() {
				if onMsg := obj.Get("onmessage"); onMsg != nil {
					if fn, ok := goja.AssertFunction(onMsg); ok {
						event := vm.NewObject()
						event.Set("data", msg)
						_, _ = fn(obj, event)
					}
				}
			})
		}

		handle, err := spawner(scriptPath, onWorkerMessage)
		if err != nil {
			panic(vm.NewGoError(err))
		}

		obj.Set("postMessage", func(call goja.FunctionCall) goja.Value {
			msg := call.Argument(0)
			handle.PostMessage(msg)
			return goja.Undefined()
		})

		obj.Set("terminate", func(call goja.FunctionCall) goja.Value {
			handle.Terminate()
			return goja.Undefined()
		})

		return obj
	})
}

func RegisterWorkerSelf(vm *goja.Runtime, postToParent func(msg goja.Value)) {
	self := vm.GlobalObject()
	vm.Set("self", self)

	self.Set("postMessage", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0)
		postToParent(msg)
		return goja.Undefined()
	})
}

// EnableWorkerAPI sets up a basic Worker constructor for standalone binaries
// Note: Worker threading in standalone binaries requires os/exec or similar,
// which is complex. For now, this is a stub that panics with a helpful message.
func EnableWorkerAPI(vm *goja.Runtime, el *eventloop.EventLoop) {
	vm.Set("Worker", func(call goja.ConstructorCall) *goja.Object {
		panic(vm.NewTypeError("Worker API is not yet supported in standalone binaries. Use 'typego run' for Worker support."))
	})
}
