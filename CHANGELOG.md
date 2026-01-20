# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.2.1] - Unreleased

### Added
- Phase 1 Foundation: `Makefile`, CI/CD Pipeline (`test.yml`), `CONTRIBUTING.md`.
- Documentation: `OPTIMIZATION.md` for performance guidelines.

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
