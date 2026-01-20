package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge"
	"github.com/repyh3/typego/bridge/core"
	"github.com/repyh3/typego/bridge/stdlib/memory"
	"github.com/repyh3/typego/bridge/stdlib/worker"
	"github.com/repyh3/typego/eventloop"

	// Import modules to trigger their init() registration
	_ "github.com/repyh3/typego/bridge/modules/fmt"
	_ "github.com/repyh3/typego/bridge/modules/net"
	_ "github.com/repyh3/typego/bridge/modules/os"
	_ "github.com/repyh3/typego/bridge/modules/sync"
)

var ErrMemoryLimitExceeded = errors.New("memory limit exceeded")

type Engine struct {
	VM            *goja.Runtime
	MemoryLimit   uint64
	EventLoop     *eventloop.EventLoop
	MemoryFactory *memory.Factory

	ctx    context.Context
	cancel context.CancelFunc
}

func NewEngine(memoryLimit uint64, mf *memory.Factory) *Engine {
	vm := goja.New()
	vm.SetMaxCallStackSize(1000)

	el := eventloop.NewEventLoop(vm)

	if mf == nil {
		mf = memory.NewFactory()
	}

	// Core modules (special dependencies)
	bridge.RegisterConsole(vm)
	// Register new stdlib modules
	memory.Register(vm, el, mf)

	// Auto-registered modules (fmt, os, http, sync)
	core.InitAll(vm, el)

	ctx, cancel := context.WithCancel(context.Background())

	eng := &Engine{
		VM:            vm,
		MemoryLimit:   memoryLimit,
		EventLoop:     el,
		MemoryFactory: mf,
		ctx:           ctx,
		cancel:        cancel,
	}

	worker.Register(vm, el, eng.SpawnWorker)

	if memoryLimit > 0 {
		eng.StartMemoryMonitor(100 * time.Millisecond)
	}

	return eng
}

func (e *Engine) Run(js string) (goja.Value, error) {
	return e.VM.RunString(js)
}

func (e *Engine) GlobalSet(name string, value interface{}) error {
	return e.VM.GlobalObject().Set(name, value)
}

func (e *Engine) BindStruct(name string, s interface{}) error {
	return core.BindStruct(e.VM, name, s)
}

func (e *Engine) Close() {
	e.cancel()
	e.EventLoop.Stop()
}

func (e *Engine) StartMemoryMonitor(interval time.Duration) {
	go func() {
		// Use a tighter interval (10ms) for critical protection
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-e.ctx.Done():
				return
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				// m.Sys is the total memory obtained from the OS.
				// This catches both JS heap and Go internal growth.
				if m.Sys > e.MemoryLimit {
					e.VM.Interrupt(ErrMemoryLimitExceeded)
					fmt.Fprintf(os.Stderr, "\n [TypeGo] CRITICAL: Memory limit reached (%d MB > %d MB). Interrupting VM...\n", m.Sys/1024/1024, e.MemoryLimit/1024/1024)
					return
				}
			}
		}
	}()
}

func (e *Engine) StartEmergencyMonitor(interval time.Duration) {
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				return
			case <-time.After(interval):
				// Placeholder for future watchdog logic
			}
		}
	}()
}
