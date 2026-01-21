package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/repyh/typego/bridge/polyfills"
	"github.com/repyh/typego/compiler"
	"github.com/repyh/typego/engine"
)

// MemoryLimit is exported for use by core commands (set in root.go)
var MemoryLimit uint64 = 128

// runInterpreter executes TypeScript directly using the embedded Goja engine.
// This is the fast path - no Go compilation required.
func runInterpreter(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Compile TypeScript to JavaScript
	res, err := compiler.Compile(absPath, nil)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Create engine with all bridge modules pre-registered
	eng := engine.NewEngine(MemoryLimit*1024*1024, nil)
	defer eng.Close()

	// Enable Node.js polyfills
	polyfills.EnableAll(eng.VM, eng.EventLoop)

	// Track execution errors
	var runErr error

	// Run on event loop
	eng.EventLoop.RunOnLoop(func() {
		_, runErr = eng.Run(res.JS)
	})

	eng.EventLoop.Start()

	if runErr != nil {
		return fmt.Errorf("runtime error: %w", runErr)
	}

	return nil
}
