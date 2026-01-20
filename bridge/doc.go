// Package bridge provides the JavaScript-to-Go binding layer for the TypeGo runtime.
//
// Bridge exposes Go functionality to the JavaScript environment running inside Goja.
// It is organized into three main areas:
//
//   - bridge/core: Low-level binding primitives and shared types.
//   - bridge/modules: Standard Go library bindings (go:fmt, go:os, go:net/http).
//   - bridge/stdlib: TypeGo-specific standard library (typego:memory, typego:worker).
//
// # Internal Registration
//
// Standard library modules in bridge/modules use a self-registration mechanism
// via init() functions and bridge/core.RegisterModule. TypeGo-specific modules
// are typically registered manually during engine initialization.
//
// # Shared Memory
//
// TypeGo provides high-performance shared memory between the main thread and workers
// via the typego:memory module.
//
// # Worker Support
//
// Background workers are supported via the typego:worker module, enabling
// multi-threaded JavaScript execution with postMessage communication.
package bridge
