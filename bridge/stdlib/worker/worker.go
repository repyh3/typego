// Package worker provides the typego:worker module for multi-threading.
package worker

import (
	"github.com/grafana/sobek"
	"github.com/repyh/typego/eventloop"
)

type Handle interface {
	PostMessage(msg sobek.Value)
	Terminate()
}

type Spawner func(scriptPath string, onMessage func(sobek.Value)) (Handle, error)

func Register(vm *sobek.Runtime, el *eventloop.EventLoop, spawner Spawner) {
	obj := vm.NewObject()

	_ = obj.Set("Worker", func(call sobek.ConstructorCall) *sobek.Object {
		scriptPath := call.Argument(0).String()

		workerObj := vm.NewObject()

		onWorkerMessage := func(msg sobek.Value) {
			el.RunOnLoop(func() {
				if onMsg := workerObj.Get("onmessage"); onMsg != nil {
					if fn, ok := sobek.AssertFunction(onMsg); ok {
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

		_ = workerObj.Set("postMessage", func(call sobek.FunctionCall) sobek.Value {
			msg := call.Argument(0)
			handle.PostMessage(msg)
			return sobek.Undefined()
		})

		_ = workerObj.Set("terminate", func(call sobek.FunctionCall) sobek.Value {
			handle.Terminate()
			return sobek.Undefined()
		})

		return workerObj
	})

	_ = vm.Set("__typego_worker__", obj)
}

func RegisterSelf(vm *sobek.Runtime, postToParent func(msg sobek.Value)) {
	self := vm.GlobalObject()
	_ = vm.Set("self", self)

	_ = self.Set("postMessage", func(call sobek.FunctionCall) sobek.Value {
		msg := call.Argument(0)
		postToParent(msg)
		return sobek.Undefined()
	})
}
