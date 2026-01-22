# TypeGo Roadmap

This document outlines the planned technical direction for **TypeGo**. Please note that these are targets; features may be implemented, deferred, or cancelled based on project needs and technical constraints.

## [v1.5.0] - Bridge Hardening & Memory
**Target**: Modernize the execution core and introduce low-level memory control.

### Core Engine & Linking
- [x] **Sobek Engine Migration**: Move from Goja to Sobek to support ES2017+ features (async/await, promises) and improve memory efficiency without using CGO.
- [ ] **Recursive Linker Upgrade**: Support for recursive type resolution, Go struct embedding, and JIT interface inspection for transparent bridging.
- [ ] **Global Intrinsics**: Support for Go-native keywords as top-level functions: `defer()`, `panic()`, and `sizeof()`.

### Low-Level Memory
- [ ] **Ref<T> API**: Implement explicit pointer management for primitives and Go structs with `.value` and `.ptr` access.
- [ ] **Pointer Handles**: Allow direct memory mutation across the bridge with safety checks and automatic finalization.
- [ ] **Alloc-Less Reflection**: Optimize `bridge/core/reflection.go` to minimize allocations during hot path execution.

### Maintenance
- [ ] **AI-Assisted Optimization**: Integration of Google Jules for daily performance audits.
- [ ] **Perf Quality Gate**: GitHub Actions that verify Jules' optimizations via automated benchmarking.

## [v1.6.0] - Native Framework & Concurrency
**Target**: Expand the built-in library with high-performance data structures and rich DX tools.

### Advanced Interop
- [ ] **TypeScript Decorators**: Initial support for `@GoType`, `@Pointer`, and lifecycle decorators for bridge configuration.
- [ ] **Shared Buffer API**: Direct memory sharing between TypeScript workers and Go for zero-copy data processing.

### Concurrency & I/O
- [ ] **Sync Primitives**: Mapping of Go’s `sync` package logic to TypeScript: `Mutex`, `WaitGroup`, and native `Channels`.
- [ ] **Streaming I/O**: TypeScript interfaces for Go’s `io.Reader` and `io.Writer` for memory-efficient processing.
- [ ] **Worker Scaling**: Improve the `typego:worker` module with auto-respawn and load balancing logic.

## [v1.7.0] - Desktop & Resource Management
**Target**: Provide tools for building lightweight desktop applications and monitoring system impact.

### Desktop APIs
- [ ] **Native Windowing**: Wails-style window management for building lightweight desktop widgets and utilities.
- [ ] **Tray & Notifications**: System-level integration for background-running TypeGo scripts.

### Resource Auditing
- [ ] **Resource Limits**: Configurable caps on memory (RAM) and CPU cycles for individual script instances.
- [ ] **Performance Telemetry**: Built-in hooks to track execution latency and bridge overhead via OpenTelemetry.

## [v1.8.0] - Developer Toolchain
**Target**: Improve the debugging and testing experience for professional TypeScript codebases.

- [ ] **Integrated Test Runner**: A native `typego test` command to execute test suites and verify Go-TypeScript interop.
- [ ] **Source Map Support**: Implementation of source maps to link runtime errors back to the original TypeScript line numbers.
- [ ] **TypeGo Studio**: A simple dashboard to monitor worker health, memory usage, and execution logs.

## [v1.9.0] - Distribution & Packaging
**Target**: Enable TypeGo logic to be shared and consumed as standard libraries.

- [ ] **NPM Library Export**: Automated packaging of TypeGo scripts into NPM modules with a Go sidecar.
- [ ] **Go Module Export**: Wrap TypeGo logic into a standard Go package for use in existing Go backends.
- [ ] **Bytecode Caching**: Persistent storage of compiled bytecode to reduce startup latency.
- [ ] **TypeGo Hub**: A curated registry for high-performance TypeGo modules.

---

*Note: Completed tasks will be marked with `[x]`. v1.4.0 is the current stable release. Deferred, cancelled, or non-listed implementations will be noted and listed above as subpoints.*