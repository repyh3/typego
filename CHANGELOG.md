# Changelog

All notable changes to this project will be documented in this file.

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
