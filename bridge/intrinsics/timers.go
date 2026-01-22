package intrinsics

import (
	"time"

	"github.com/grafana/sobek"
)

// EnableTimers injects setTimeout and setInterval globals
func (r *Registry) EnableTimers() {
	_ = r.vm.Set("setTimeout", func(call sobek.FunctionCall) sobek.Value {
		fn, _ := sobek.AssertFunction(call.Argument(0))
		ms := call.Argument(1).ToInteger()
		r.el.WGAdd(1)
		go func() {
			time.Sleep(time.Duration(ms) * time.Millisecond)
			r.el.RunOnLoop(func() {
				_, _ = fn(sobek.Undefined())
				r.el.WGDone()
			})
		}()
		return sobek.Undefined()
	})

	_ = r.vm.Set("setInterval", func(call sobek.FunctionCall) sobek.Value {
		fn, _ := sobek.AssertFunction(call.Argument(0))
		ms := call.Argument(1).ToInteger()

		stop := make(chan struct{})

		go func() {
			ticker := time.NewTicker(time.Duration(ms) * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					r.el.RunOnLoop(func() {
						_, _ = fn(sobek.Undefined())
					})
				case <-stop:
					return
				}
			}
		}()

		// Return a simple ID object for clearInterval
		id := r.vm.NewObject()
		_ = id.Set("__stop__", stop)
		return id
	})

	_ = r.vm.Set("clearInterval", func(call sobek.FunctionCall) sobek.Value {
		obj := call.Argument(0).ToObject(r.vm)
		if obj != nil {
			if ch := obj.Get("__stop__"); ch != nil {
				if stop, ok := ch.Export().(chan struct{}); ok {
					close(stop)
				}
			}
		}
		return sobek.Undefined()
	})
}
