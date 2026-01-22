# Changelog

All notable changes to this project will be documented in this file.

## [v1.5.0] - 2026-01-22

### Core Engine Upgrade
- **Sobek Migration**: Successfully migrated from Goja to Sobek. This enables modern ES2017+ features (async/await, promises) and improves performance for complex TypeScript workloads without CGO dependencies.
- **Thread-Safe Runtime**: Introduced a global `VMLock` and a non-blocking `select` loop that yields this lock, allowing Go goroutines to safely interact with the single-threaded JavaScript environment.

### Global Intrinsics (The "Low-Level" Milestone)
- **`panic()` & `recover()`**: Implemented Go-native error handling. `panic()` stops execution with a Go-formatted error, and `recover()` allows safe capture within deferred functions.
- **`defer()`**: Added a compile-time transformation that schedules cleanup tasks to run at the end of a function, even if a panic occurs.
- **`sizeof()`**: Implemented a deep memory inspection intrinsic with circular reference detection and meta-caching for high-performance sizing.
- **`go()`**: Direct bridging to Go goroutines. Launch any JS function effectively in the background.
- **`makeChan()` & `select()`**: Native Go channels and multiplexing for robust concurrent communication between JS tasks.
- **`iota`**: Support for auto-incrementing compile-time constants.

### Memory & Pointers
- **Pointer Intrinsics**: Introduced `ref()` and `deref()` for explicit variable tracking and reference management.
- **Slice Primitives**: Added `make()`, `cap()`, and `copy()` to provide Go-like performance for TypedArray manipulation.

### DX & CLI
- **New Examples**: Revamped the `examples/` directory with detailed showcases for concurrency, memory management, and advanced control flow.
- **IntelliSense Expansion**: Full TypeScript definitions for all new intrinsics integrated into the automatic type generation pipeline.

## [v1.4.0] - 2026-01-21

### TypeGo Package Ecosystem
- **New Project Architecture**: Unified project structure with `typego.modules.json` manifest and `.typego/` hidden workspace.
- **Dependency Management CLI**: A complete suite of commands for managing Go dependencies:
    - `typego add <pkg>[@version]`: Adds Go modules to the project and triggers auto-type generation.
    - `typego remove <pkg>`: Cleanly removes dependencies and updates types.
    - `typego list`: Visualizes the dependency tree and current versions.
    - `typego update/outdated`: Tools for keeping Go modules up to date.
    - `typego install`: Manually triggers JIT build, dependency resolution, and type syncing.
- **Reproducible Builds**: Introduced `typego.lock` to pin resolved Go module versions across different environments.
- **JIT Infrastructure**: Automated workspace management with `go mod tidy` integration and checksum-based stale detection for faster re-installs.

### Smarter Type Linker
- **Recursive Type Discovery**: Automatically discovers and generates TypeScript interfaces for structs/interfaces returned by functions or nested in fields.
- **Go Struct Embedding**: Mirrored struct embedding in TypeScript using the `extends` keyword (e.g., `Engine extends RouterGroup`).
- **Local Reference Resolution**: Correctly resolves and links internal struct references within the same package.
- **Convention Alignment**: Standard library types (e.g., `go:net/http`) now use Go-style PascalCase properties (`Method`, `URL`, `Header`) for consistency with auto-generated external types.

### CLI & DX Improvements
- **CLI Restructuring**: Refactored the internal CLI architecture for better maintainability and extensibility.
- **Automated Workflows**: Automatic type regeneration on `add` and `install` commands.
- **Enhanced Scaffolding**: `typego init` now generates a production-ready setup including `.gitignore` and pre-configured `tsconfig.json`.

### Changed
- Replaced manual JIT build logic with a unified `installer` package.
- Updated `typego run` to check for module staleness before execution.

## [v1.3.1] - 2026-01-21

### Fixed
- Fixed internal syntax error in `typego:worker` type definitions (`worker.d.ts`).
- Decentralized `go:crypto` types to ensure correct linking and type generation.

## [v1.3.0] - 2026-01-20

### Added
- **Bindings**: `go:encoding/json`, `go:os` (Getenv, Exit, etc.), `go:crypto/*` (Sha256, Hmac, Rand).
- **HTTP Server**: `ListenAndServe` with JS callback, `Request`/`Response` objects, and `Post`.
- **Developer Experience**: `typego dev` command with hot-reload and colored output.
- **Cross-Compilation**: `typego build --target` support for Linux, macOS, and Windows.
- **Runtime**: Context support and graceful shutdown via `server.close()`.

## [v1.2.0] - 2026-01-20

### Added
- **Interpreter Mode**: Fast execution (~0.2s startup) using `goja` directly.
- **Modular Stdlib**: Refactored `bridge` into `core`, `modules`, and `stdlib`.
- **Native Modules**: `typego:memory` and `typego:worker` available natively.
- **CLI**: `typego run` now defaults to interpreter mode. Use `--compile` for standalone binary.

### Changed
- Refactored `cmd/typego/run.go` to support dual execution modes.
- Updated `README.md` with concise styles and tables.
- Removed legacy `examples/cluster` directory.

### Fixed
- Fixed internal module hyper-linking in standalone builds.
- Improved dependency resolution and tidy steps in compilation.

## [v1.0.0] - Initial Release

- Core TypeGo Runtime
- Goja Integration
- Basic Examples
