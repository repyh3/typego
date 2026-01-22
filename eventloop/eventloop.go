package eventloop

import (
	"context"
	"sync"

	"github.com/grafana/sobek"
)

type RejectionHandler func(err error)

type EventLoop struct {
	VM       *sobek.Runtime
	jobQueue chan func()
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
	autoStop bool

	ctx    context.Context
	cancel context.CancelFunc

	OnUnhandledRejection RejectionHandler
}

func NewEventLoop(vm *sobek.Runtime) *EventLoop {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventLoop{
		VM:       vm,
		jobQueue: make(chan func(), 100),
		stopChan: make(chan struct{}),
		autoStop: true,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (el *EventLoop) SetAutoStop(enable bool) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.autoStop = enable
}

// Start runs the event loop. It blocks until the loop is stopped or all tasks are done.
func (el *EventLoop) Start() {
	el.mu.Lock()
	if el.running {
		el.mu.Unlock()
		return
	}
	el.running = true
	shouldAutoStop := el.autoStop
	el.mu.Unlock()

	// Shutdown when no more tasks are pending
	if shouldAutoStop {
		go func() {
			el.wg.Wait()
			el.Stop()
		}()
	}

	for {
		select {
		case job := <-el.jobQueue:
			job()
			el.wg.Done()
		case <-el.stopChan:
			return
		}
	}
}

// RunOnLoop schedules a function to run on the JS thread. Safe for concurrent use.
func (el *EventLoop) RunOnLoop(f func()) {
	el.wg.Add(1)
	el.jobQueue <- f
}

func (el *EventLoop) Stop() {
	el.mu.Lock()
	defer el.mu.Unlock()
	if !el.running {
		return
	}
	el.cancel()
	close(el.stopChan)
	el.running = false
}

func (el *EventLoop) Context() context.Context {
	return el.ctx
}

func (el *EventLoop) Shutdown(timeout context.Context) error {
	done := make(chan struct{})
	go func() {
		el.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		el.Stop()
		return nil
	case <-timeout.Done():
		el.Stop()
		return timeout.Err()
	}
}

func (el *EventLoop) WGAdd(n int) {
	el.wg.Add(n)
}

func (el *EventLoop) WGDone() {
	el.wg.Done()
}

func (el *EventLoop) CreatePromise() (promise *sobek.Object, resolve func(interface{}), reject func(interface{})) {
	p, res, rej := el.VM.NewPromise()

	// Keep the loop alive until the promise is settled
	el.wg.Add(1)

	resolve = func(v interface{}) {
		el.RunOnLoop(func() {
			_ = res(v)
			el.wg.Done()
		})
	}

	reject = func(v interface{}) {
		el.RunOnLoop(func() {
			_ = rej(v)
			el.wg.Done()
		})
	}

	return el.VM.ToValue(p).ToObject(el.VM), resolve, reject
}
