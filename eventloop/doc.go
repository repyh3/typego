// The event loop ensures that all Goja runtime operations occur on a single
// goroutine, preventing race conditions. It also manages async operations
// through a wait group mechanism.
//
// # Basic Usage
//
//	el := eventloop.NewEventLoop(vm)
//
//	el.RunOnLoop(func() {
//	    vm.RunString(`console.log("Hello")`)
//	})
//
//	el.Start() // Blocks until all work is done
//
// # Async Operations
//
// When starting an async operation (HTTP request, timer, etc.), use WGAdd to
// indicate pending work and WGDone when complete:
//
//	el.WGAdd(1)
//	go func() {
//	    time.Sleep(time.Second)
//	    el.RunOnLoop(func() {
//	        vm.RunString(`console.log("Delayed")`)
//	        el.WGDone()
//	    })
//	}()
//
// The event loop automatically stops when all pending operations complete.
//
// # Promises
//
// CreatePromise returns a JavaScript Promise along with resolve/reject functions
// that properly integrate with the event loop:
//
//	promise, resolve, reject := el.CreatePromise()
//	go func() {
//	    result, err := doAsyncWork()
//	    if err != nil {
//	        reject(err)
//	    } else {
//	        resolve(result)
//	    }
//	}()
//	return promise
package eventloop
