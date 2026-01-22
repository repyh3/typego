package engine

import (
	"fmt"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/bridge/stdlib/worker"
	"github.com/repyh/typego/compiler"
)

type WorkerInstance struct {
	vm     *sobek.Runtime
	engine *Engine
	inbox  chan interface{}
	stop   chan struct{}
}

func (w *WorkerInstance) PostMessage(msg sobek.Value) {
	data := msg.Export()
	w.inbox <- data
}

func (w *WorkerInstance) Terminate() {
	close(w.stop)
}

func (e *Engine) SpawnWorker(scriptPath string, onMessage func(sobek.Value)) (worker.Handle, error) {
	res, err := compiler.Compile(scriptPath, nil)
	if err != nil {
		return nil, fmt.Errorf("compile error: %w", err)
	}

	workerEng := NewEngine(e.MemoryLimit, e.MemoryFactory)

	inbox := make(chan interface{}, 100)
	stop := make(chan struct{})

	w := &WorkerInstance{
		vm:     workerEng.VM,
		engine: workerEng,
		inbox:  inbox,
		stop:   stop,
	}

	worker.RegisterSelf(workerEng.VM, func(msg sobek.Value) {
		data := msg.Export()
		e.EventLoop.RunOnLoop(func() {
			val := e.VM.ToValue(data)
			onMessage(val)
		})
	})

	go func() {
		_, err := workerEng.Run(res.JS)
		if err != nil {
			fmt.Printf("Worker Error [%s]: %v\n", scriptPath, err)
		}

		workerEng.EventLoop.WGAdd(1)
		workerEng.EventLoop.Start()
	}()

	go func() {
		for {
			select {
			case msg := <-inbox:
				workerEng.EventLoop.RunOnLoop(func() {
					val := workerEng.VM.ToValue(msg)
					if onMsg := workerEng.VM.GlobalObject().Get("onmessage"); onMsg != nil {
						if fn, ok := sobek.AssertFunction(onMsg); ok {
							event := workerEng.VM.NewObject()
							_ = event.Set("data", val)
							_, _ = fn(workerEng.VM.GlobalObject(), event)
						}
					}
				})
			case <-stop:
				workerEng.EventLoop.Stop()
				return
			}
		}
	}()

	return w, nil
}
