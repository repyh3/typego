package integration

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/engine"
)

// TestHarness bundles common test infrastructure
type TestHarness struct {
	Engine *engine.Engine
	Ctx    context.Context
}

// NewHarness creates a fresh engine instance for testing
func NewHarness(t *testing.T) *TestHarness {
	// Initialize with standard memory limit (64MB)
	eng := engine.NewEngine(64*1024*1024, nil)

	// Ensure we don't leak engines if test panics/fails
	t.Cleanup(func() {
		eng.Close()
	})

	// Disable auto-stop so we can reuse the loop
	eng.EventLoop.SetAutoStop(false)
	go eng.EventLoop.Start()

	return &TestHarness{
		Engine: eng,
		Ctx:    context.Background(),
	}
}

func (h *TestHarness) Run(t *testing.T, code string) {
	done := make(chan error, 1)

	h.Engine.EventLoop.RunOnLoop(func() {
		// Wrap code in async IIFE to allow await (simulating module)
		// Actually, Goja RunString doesn't support top-level await even with IIFE unless we return the promise.
		// We'll wrap in `async function main() { ... }; main()`
		wrappedCode := "(async () => { " + code + " })()"

		val, err := h.Engine.Run(wrappedCode)
		if err != nil {
			done <- err
			return
		}

		// Handle Promise
		if val != nil && !sobek.IsUndefined(val) && !sobek.IsNull(val) {
			if obj := val.ToObject(h.Engine.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !sobek.IsUndefined(then) {
					if _, ok := sobek.AssertFunction(then); ok {
						h.Engine.EventLoop.WGAdd(1)

						onDone := h.Engine.VM.ToValue(func(sobek.FunctionCall) sobek.Value {
							h.Engine.EventLoop.WGDone()
							done <- nil // Success
							return sobek.Undefined()
						})

						onErr := h.Engine.VM.ToValue(func(call sobek.FunctionCall) sobek.Value {
							h.Engine.EventLoop.WGDone()
							// call.Argument(0) is the error
							errVal := call.Argument(0)
							done <- jsError{errVal}
							return sobek.Undefined()
						})

						thenFn, _ := sobek.AssertFunction(then)
						_, _ = thenFn(val, onDone, onErr)
						return
					}
				}
			}
		}
		// Not a promise, done immediately
		done <- nil
	})

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Runtime Error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for script execution")
	}
}

type jsError struct {
	val sobek.Value
}

func (e jsError) Error() string { return e.val.String() }
