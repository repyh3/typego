<div align="center">

# TypeGo

typego is an embedded TypeScript runtime for Go. It lets you script Go applications with TS without the overhead of Node.js or the boilerplate of manual FFI bindings.

[Getting Started](#getting-started) • [Features](#features) • [Examples](#examples) • [Optimization](OPTIMIZATION.md) • [Contributing](CONTRIBUTING.md) • [License](#license)

</div>

> [!NOTE]
> **Project Status**: TypeGo is under active development. However, please note that **maintenance is limited** as I am balancing this project with my university commitments. Issues and PRs are welcome but may see delayed responses.

Unlike typical runtimes that communicate over IPC or JSON-RPC, typego runs a JS engine (Goja) directly inside your Go process. You can import Go packages as if they were native TS modules using the go: prefix, allowing for zero-copy data sharing and direct access to Go’s standard library.

## Features

- **Direct Go Integration**: Import any Go package as a native TS module (`go:fmt`, `go:github.com/gin-gonic/gin`).
- **Smarter Type Linker**: Automatic, recursive type generation for Go structs, interfaces, and methods. Supports struct embedding (`extends`) and nested type resolution.
- **True Parallelism**: Goroutine-based workers with zero-copy shared memory (`typego:memory`).
- **Modern Package Ecosystem**: Built-in CLI for managing Go dependencies with `typego.modules.json` and `typego.lock`.
- **Fast Developer Loop**: Hot-reloading dev server and ~0.2s interpreter startup. Compiles to single-binary with `--compile`.

## Tech Stack

- **Language**: Go 1.21+, TypeScript
- **JS Engine**: Goja (pure Go, no CGO)
- **Bundler**: esbuild
- **CLI**: Cobra

## Getting Started

### Installation

```bash
go install github.com/repyh/typego/cmd/typego@latest
```

### Quick Start

```bash
typego init myapp
cd myapp
typego run src/index.ts
```

## Project Structure

A standard TypeGo project consists of the following structure:

- `src/`: Directory for your TypeScript source files.
- `typego.modules.json`: Dependency manifest for Go packages.
- `typego.lock`: **[New]** Auto-generated lockfile for reproducible Go module versions.
- `.typego/`: **[Internal]** Managed workspace for build artifacts and cached TypeScript types (`go.d.ts`).
- `package.json`: Standard Node.js manifest for NPM dependencies (handled via esbuild).

### Commands

| Command | Description |
|---------|-------------|
| `typego run <file>` | Execute TypeScript (fast interpreter mode) |
| `typego dev <file>` | Development server with hot-reload |
| `typego build <file>` | Build standalone executable |
| `typego init [name]` | Scaffold new project (`--npm` for Node interop) |
| `typego types` | Generate `.d.ts` for Go imports |
| `typego add <pkg>` | Add a Go module dependency |
| `typego remove <pkg>` | Remove a Go module dependency |
| `typego list` | List configured Go dependencies |
| `typego update` | Update Go modules to latest versions |
| `typego outdated` | Check for newer Go module versions |
| `typego install` | Manually trigger JIT build/dependency resolution |
| `typego clean` | Reset build cache and temporary workspace |

### Package Management

TypeGo uses a `typego.modules.json` file to manage Go dependencies. This allows you to use any Go package in your TypeScript code.

```bash
# Add a Go package. Ecosystem will automatically resolve versions and sync types.
typego add github.com/gin-gonic/gin
```
```typescript
// src/index.ts
import { Default } from "go:github.com/gin-gonic/gin";

const app = Default();
app.GET("/ping", (c) => c.JSON(200, { message: "pong" }));
app.Run();
```

TypeGo automatically manages a `.typego/` workspace, handling `go mod tidy`, JIT compilation, and TypeScript definition syncing behind the scenes.

## Examples

### Go Imports
```typescript
import { Println, Printf } from "go:fmt";
import { Sleep } from "go:sync";

Println("Hello from Go!");
await Sleep(1000);
Printf("Done after %dms\n", 1000);
```

### Workers & Shared Memory
```typescript
import { makeShared } from "typego:memory";
import { Worker } from "typego:worker";

const shared = makeShared("buffer", 1024);
const worker = new Worker("worker.ts");
worker.postMessage({ buffer: shared });
```

### Concurrency
```typescript
import { Spawn, Sleep } from "go:sync";

Spawn(async () => {
  for (let i = 0; i < 5; i++) {
    await Sleep(500);
    console.log(`Heartbeat ${i}`);
  }
});
```

## Applications

| Field | Can Do | Cannot Do |
|-------|--------|-----------|
| **Networking** | HTTP/WebSocket/TCP servers, proxies | Kernel networking, eBPF |
| **Backend** | REST APIs, microservices, job queues | Native DB drivers (no CGO) |
| **CLI Tools** | Build tools, code generators, automation | Native GUI |
| **Data Processing** | JSON/CSV, log aggregation, parallel ETL | GPU acceleration |
| **DevOps** | Health checks, log shippers, K8s clients | Container runtimes |
| **Real-Time** | Chat servers, game servers, notifications | Hard real-time |

## Limitations

| Limitation | Impact |
|------------|--------|
| No JIT compilation | ~10x slower raw JS than V8-based runtimes |
| No Web APIs | No DOM, fetch, localStorage |
| No CGO | Can't call C libraries |
| Explicit shared memory | Must use `makeShared()` |

## When to Use

| ✅ Good Fit | ❌ Use Alternatives |
|------------|---------------------|
| Backend needing true parallelism | Browser/frontend apps |
| CLI tools with single binary | Data science (use Python) |
| TypeScript team wanting Go perf | Systems programming (use Rust/Go) |
| Serverless/edge (small binary) | Mobile apps |

## Performance

TypeGo is optimized for high-throughput I/O and true parallelism. For a detailed breakdown of the execution model, interop overhead, and optimization strategies, see **[OPTIMIZATION.md](OPTIMIZATION.md)**.


## Runtime Comparison

| Feature | TypeGo | Node.js | Deno | Bun |
|---------|:------:|:-------:|:----:|:---:|
| TypeScript native | ✅ | ⚠️ | ✅ | ✅ |
| True parallelism | ✅ | ❌ | ❌ | ❌ |
| Single binary | ✅ | ❌ | ✅ | ❌ |
| Shared memory | ✅ | ⚠️ | ⚠️ | ⚠️ |
| NPM ecosystem | ⚠️ | ✅ | ⚠️ | ✅ |

## Development

### Prerequisites

- Go 1.21+
- Node.js 18+ (for NPM packages)

### Building from Source

```bash
git clone https://github.com/repyh/typego.git
cd typego
go build -o typego.exe ./cmd/typego
```

### Running Examples

```bash
./typego run examples/01-hello-world.ts
./typego run examples/02-go-concurrency.ts
./typego run examples/10-typego-stdlib.ts
```

## License

MIT
