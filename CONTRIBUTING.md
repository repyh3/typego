# Contributing to TypeGo

Thank you for your interest in TypeGo! We adhere to **Professional Open Source Standards**. Please review this guide to ensure a smooth collaboration.

## Quick Reference

Common commands for development:

| Action | Command | Description |
| :--- | :--- | :--- |
| **Setup** | `go mod download` | Install Go dependencies. |
| **Dev Mode** | `typego watch <file>` | Run script with hot-reload. |
| **Test (Unit)** | `go test ./engine/...` | Test core runtime only. |
| **Test (Full)** | `go test -race ./...` | **Required.** Run all tests with race detection. |
| **Build CLI** | `go build ./cmd/typego` | Compile the TypeGo CLI binary. |
| **Lint** | `golangci-lint run` | Verify code quality. |

---

## Development Workflow

We follow a strict **Feature Branch** workflow. Direct commits to `main` are blocked.

| Stage | Command / Action | Output / Goal |
| :--- | :--- | :--- |
| **1. Branch** | `git checkout -b <type>/<name>` | Create a focused branch (e.g., `feat/http-client` or `fix/mem-leak`). |
| **2. Code** | *Write idiomatic Go code* | Follow `Effective Go`. Keep it simple and performant. |
| **3. Verify** | `go test -race -v ./...` | Ensure **0** race conditions and **100%** pass rate. |
| **4. E2E** | `go test ./tests/e2e/...` | Verify the CLI builds and runs correctly. |
| **5. Commit** | `git commit -m "feat: add support for..."` | Use **Conventional Commits** (see below). |
| **6. Push** | `git push -u origin <branch>` | Push to your fork. |

---

## Pull Request Guidelines

To ensure rapid review and merging, please adhere to these rules:

### 1. Branch Naming Strategy
| Prefix | Use Case | Example |
| :--- | :--- | :--- |
| `feat/` | New features or enhancements | `feat/async-io` |
| `fix/` | Bug fixes | `fix/race-condition` |
| `docs/` | Documentation only | `docs/readme-update` |
| `test/` | Adding missing tests | `test/http-integration` |
| `perf/` | Performance optimizations | `perf/cache-layer` |

### 2. The PR Checklist
Before submitting, verify the following:
- [ ] **Tests Passed**: `go test -race ./...` is green.
- [ ] **New Tests Added**: If fixing a bug, include a reproduction test case.
- [ ] **Lint Free**: No linter warnings.
- [ ] **No Dead Code**: Remove all `fmt.Println` debug statements.

### 3. Commit Message Convention
We use [Conventional Commits](https://www.conventionalcommits.org/).
- **Format**: `<type>(<scope>): <description>`
- **Example**: `feat(engine): implement strict memory limits`
- **Example**: `fix(bridge): resolve race condition in EventLoop`

---

## Project Structure

| Directory | Purpose |
| :--- | :--- |
| `cmd/typego` | The CLI entry point (`main`, `run`, `build`). |
| `engine/` | Core Goja runtime wrapper and EventLoop. |
| `bridge/` | Bindings between Go and JS. |
| `compiler/` | ESBuild wrapper containing the **Cache** logic. |
| `tests/integration/` | Verifies `bridge` modules function correctly in JS. |
| `tests/e2e/` | Verifies the CLI binary compiles user apps correctly. |

---

## Need Help?
If you're unsure about an implementation detail:
1. Check `TESTING_STRATEGY.md` for architectural context.
2. Open an Issue with the label `question`.
