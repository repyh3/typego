package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/repyh/typego/compiler"
	"github.com/repyh/typego/engine"
)

var MemoryLimit uint64 = 128

// runInterpreter executes TypeScript directly using the embedded Goja engine.
// This is the fast path - no Go compilation required.
func runInterpreter(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	res, err := compiler.Compile(absPath, nil)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	eng := engine.NewEngine(MemoryLimit*1024*1024, nil)
	defer eng.Close()

	var runErr error

	eng.EventLoop.RunOnLoop(func() {
		_, runErr = eng.Run(res.JS)
	})

	eng.EventLoop.Start()

	if runErr != nil {
		return fmt.Errorf("runtime error: %w", runErr)
	}

	return nil
}
