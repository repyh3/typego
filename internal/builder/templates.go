package builder

// ShimTemplate is the Go source template used to generate standalone TypeGo binaries.
// It wraps the bundled JavaScript in a Goja runtime with all bridge modules enabled.
//
// Template placeholders:
//   - %[1]s: Additional import statements for hyper-linked Go packages
//   - %[2]s: The bundled JavaScript code (quoted string)
//   - %[3]s: Hyper-linker binding code (generated shims)
//   - %[4]d: Memory limit in bytes
const ShimTemplate = `package main

import (
	"fmt"
	"os"

	%[1]s

	"github.com/grafana/sobek"
	"github.com/repyh/typego/engine"
)

const jsBundle = %[2]s

type NativeTools struct {
	StartTime string
}

func (n *NativeTools) GetRuntimeInfo() string {
	return "TypeGo Standalone v1.0"
}

func main() {
	// Use the unified engine which handles ALL module registration
	eng := engine.NewEngine(%[4]d, nil)

	// Initialize Native Tools
	tools := &NativeTools{StartTime: "2026-01-20"}
	_ = eng.BindStruct("native", tools)

	// Hyper-Linker Bindings (Generated)
	%[3]s

	// Run on EventLoop
	eng.EventLoop.RunOnLoop(func() {
		val, err := eng.Run(jsBundle)
		if err != nil {
			fmt.Printf("Runtime Error: %%v\n", err)
			os.Exit(1)
		}

		// Handle Top-Level Async (Promises)
		if val != nil && !sobek.IsUndefined(val) && !sobek.IsNull(val) {
			if obj := val.ToObject(eng.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !sobek.IsUndefined(then) {
					if _, ok := sobek.AssertFunction(then); ok {
						eng.EventLoop.WGAdd(1)
						done := eng.VM.ToValue(func(sobek.FunctionCall) sobek.Value {
							eng.EventLoop.WGDone()
							return sobek.Undefined()
						})
						thenFn, _ := sobek.AssertFunction(then)
						_, _ = thenFn(val, done, done)
					}
				}
			}
		}
	})

	eng.EventLoop.Start()
}
`
