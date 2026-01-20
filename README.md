<div align="center">

# TypeGo

typego is an embedded TypeScript runtime for Go. It lets you script Go applications with TS without the overhead of Node.js or the boilerplate of manual FFI bindings.

[Getting Started](#getting-started) • [Features](#features) • [Examples](#examples) • [Optimization](OPTIMIZATION.md) • [Contributing](CONTRIBUTING.md) • [License](#license)

</div>

> [!NOTE]
> **Project Status**: TypeGo is under active development. However, please note that **maintenance is limited** as I am balancing this project with my university commitments. Issues and PRs are welcome but may see delayed responses.

Unlike typical runtimes that communicate over IPC or JSON-RPC, typego runs a JS engine (Goja) directly inside your Go process. You can import Go packages as if they were native TS modules using the go: prefix, allowing for zero-copy data sharing and direct access to Go’s standard library.

## Features

- **True Parallelism**: Goroutine-based workers that run in parallel, not just concurrently.
- **Single Binary**: Compiles TypeScript + Go runtime into one executable. No `node_modules`.
- **Shared Memory**: Zero-copy ArrayBuffers between workers with mutex protection.
- **Go Imports**: Use Go's standard library directly (`go:fmt`, `go:net/http`, `go:sync`).
- **NPM Compatible**: Standard NPM packages work via esbuild bundling.
- **Fast Interpreter**: Default mode executes in ~0.2s. Use `--compile` for standalone binaries.

## Tech Stack

- **Language**: Go 1.21+, TypeScript
- **JS Engine**: Goja (pure Go, no CGO)
- **Bundler**: esbuild
- **CLI**: Cobra

## Getting Started

### Installation

```bash
go install github.com/repyh3/typego/cmd/typego@latest
```

### Quick Start

```bash
typego init myapp
cd myapp
typego run src/index.ts
```

### Commands

| Command | Description |
|---------|-------------|
| `typego run <file>` | Execute TypeScript (fast interpreter mode) |
| `typego dev <file>` | Development server with hot-reload |
| `typego run --compile <file>` | Compile and run as standalone binary |
| `typego build <file> -o <out>` | Build standalone executable |
| `typego build <file> --target` | Cross-compile (e.g. linux-amd64) |
| `typego types` | Generate `.d.ts` for Go imports |
| `typego init <name>` | Scaffold new project |

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
git clone https://github.com/repyh3/typego.git
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
