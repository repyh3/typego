package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/bridge/intrinsics"
	"github.com/repyh/typego/bridge/stdlib/memory"
	"github.com/repyh/typego/bridge/stdlib/worker"
	"github.com/repyh/typego/eventloop"

	_ "github.com/repyh/typego/bridge/intrinsics"
	_ "github.com/repyh/typego/bridge/modules/crypto"
	_ "github.com/repyh/typego/bridge/modules/fmt"
	_ "github.com/repyh/typego/bridge/modules/json"
	_ "github.com/repyh/typego/bridge/modules/net"
	_ "github.com/repyh/typego/bridge/modules/os"
	_ "github.com/repyh/typego/bridge/modules/sync"
)

var ErrMemoryLimitExceeded = errors.New("memory limit exceeded")

// This is used by JIT binaries to register custom modules.
type GlobalEngineHook func(eng *Engine)

var GlobalHooks []GlobalEngineHook

func AddGlobalHook(hook GlobalEngineHook) {
	GlobalHooks = append(GlobalHooks, hook)
}

type ErrorHandler func(err error, stack string)

type Engine struct {
	VM            *sobek.Runtime
	MemoryLimit   uint64
	EventLoop     *eventloop.EventLoop
	MemoryFactory *memory.Factory
	Intrinsics    *intrinsics.Registry

	// OnError is called when an unhandled error occurs in the engine

	OnError ErrorHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func (e *Engine) WrapError(recovered interface{}) error {
	switch v := recovered.(type) {
	case *sobek.Exception:
		return v
	case error:
		return fmt.Errorf("runtime error: %w", v)
	case string:
		return fmt.Errorf("runtime error: %s", v)
	default:
		return fmt.Errorf("runtime error: %v", v)
	}
}

func NewEngine(memoryLimit uint64, mf *memory.Factory) *Engine {
	vm := sobek.New()
	vm.SetMaxCallStackSize(1000)

	el := eventloop.NewEventLoop(vm)

	if mf == nil {
		mf = memory.NewFactory()
	}

	core.RegisterConsole(vm)
	core.RegisterGlobals(vm)

	memory.Register(vm, el, mf)

	core.InitAll(vm, el)

	ctx, cancel := context.WithCancel(context.Background())

	// Enable Global Intrinsics (panic, sizeof, typego.scope)
	intrinsicsReg := intrinsics.Enable(vm, el)

	eng := &Engine{
		VM:            vm,
		MemoryLimit:   memoryLimit,
		EventLoop:     el,
		MemoryFactory: mf,
		Intrinsics:    intrinsicsReg,
		ctx:           ctx,
		cancel:        cancel,
	}

	worker.Register(vm, el, eng.SpawnWorker)

	if memoryLimit > 0 {
		eng.StartMemoryMonitor(100 * time.Millisecond)
	}

	// Apply global hooks
	for _, hook := range GlobalHooks {
		hook(eng)
	}

	return eng
}

func (e *Engine) Run(js string) (sobek.Value, error) {
	e.Intrinsics.VMLock.Lock()
	defer e.Intrinsics.VMLock.Unlock()
	return e.VM.RunString(js)
}

// RunSafe executes JS code with panic recovery. If a panic occurs, it is
// converted to an error and passed to OnError if set.
func (e *Engine) RunSafe(js string) (result sobek.Value, err error) {
	e.Intrinsics.VMLock.Lock()
	defer e.Intrinsics.VMLock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			err = e.WrapError(r)
			if e.OnError != nil {
				e.OnError(err, e.getStack())
			}
		}
	}()
	return e.VM.RunString(js)
}

func (e *Engine) getStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func (e *Engine) Context() context.Context {
	return e.ctx
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
