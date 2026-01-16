package engine

import (
	"time"

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge"
	"github.com/repyh3/typego/eventloop"
)

type Engine struct {
	VM            *goja.Runtime
	MemoryLimit   uint64
	EventLoop     *eventloop.EventLoop
	MemoryFactory *bridge.MemoryFactory
}

func NewEngine(memoryLimit uint64, mf *bridge.MemoryFactory) *Engine {
	vm := goja.New()
	vm.SetMaxCallStackSize(1000)

	el := eventloop.NewEventLoop(vm)

	if mf == nil {
		mf = bridge.NewMemoryFactory()
	}

	bridge.RegisterConsole(vm)
	bridge.RegisterMemory(vm, mf, el)
	bridge.RegisterFmt(vm)
	bridge.RegisterOS(vm)
	bridge.RegisterHTTP(vm, el)
	bridge.RegisterSync(vm, el)
	bridge.RegisterMemoryFactory(vm, mf, el)

	eng := &Engine{
		VM:            vm,
		MemoryLimit:   memoryLimit,
		EventLoop:     el,
		MemoryFactory: mf,
	}

	bridge.RegisterWorker(vm, el, eng.SpawnWorker)

	return eng
}

func (e *Engine) Run(js string) (goja.Value, error) {
	return e.VM.RunString(js)
}

func (e *Engine) GlobalSet(name string, value interface{}) error {
	return e.VM.GlobalObject().Set(name, value)
}

func (e *Engine) BindStruct(name string, s interface{}) error {
	return bridge.BindStruct(e.VM, name, s)
}

func (e *Engine) StartEmergencyMonitor(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
		}
	}()
}
