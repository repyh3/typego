<div align="center">

# TypeGo

typego is an embedded TypeScript runtime for Go. It lets you script Go applications with TS without the overhead of Node.js or the boilerplate of manual FFI bindings.

[Getting Started](#getting-started) • [Features](#features) • [Examples](#examples) • [Optimization](OPTIMIZATION.md) • [Contributing](CONTRIBUTING.md) • [License](#license)

</div>

> [!NOTE]
> **Project Status**: TypeGo is under active development. However, please note that **maintenance is limited** as I am balancing this project with my university commitments. Issues and PRs are welcome but may see delayed responses.

Unlike typical runtimes that communicate over IPC or JSON-RPC, typego runs a JS engine (Sobek) directly inside your Go process. You can import Go packages as if they were native TS modules using the go: prefix, allowing for zero-copy data sharing and direct access to Go’s standard library.

## Index

- [Overview](#overview)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Usage](#usage)
- [Language Reference](#language-reference)
  - [Imports](#imports)
  - [Intrinsics](#intrinsics)
    - [Process & Environment](#process--environment)
    - [Encoding](#encoding)
    - [IO Utilities](#io-utilities)
  - [Standard Library](#standard-library)
    - [go:fmt](#gofmt)
    - [go:net/http](#gonethttp)
    - [go:os](#goos)
    - [go:sync](#gosync)
    - [typego:memory](#typegomemory)
    - [typego:worker](#typegoworker)
  - [Concurrency](#concurrency)
    - [Goroutines](#goroutines)
    - [Channels](#channels)
    - [Select](#select)
  - [Memory Management](#memory-management)
    - [Pointers & Refs](#pointers--refs)
    - [Shared Memory](#shared-memory)
    - [Defer](#defer)
- [Tooling](#tooling)
  - [CLI Reference](#cli-reference)
  - [Package Management](#package-management)
- [Ecosystem](#ecosystem)
  - [Performance](#performance)
  - [Runtime Comparison](#runtime-comparison)
- [Cookbook](#cookbook)
  - [HTTP Server with Middleware](#http-server-with-middleware)
  - [Worker Pool](#worker-pool)
  - [File Processing Pipeline](#file-processing-pipeline)
  - [External Go Modules](#external-go-modules)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [License](#license)

---

## Overview

TypeGo bridges the gap between Go's raw performance and TypeScript's developer experience. It is not just a JS runtime; it is a **hybrid runtime** where TypeScript code compiles JIT into the Go process, sharing memory and goroutines.

### Features

- **Direct Go Integration**: Import any Go package as a native TS module (`go:fmt`, `go:github.com/gin-gonic/gin`).
- **Standard Library Intrinsics**: Direct access to Go-native keywords and types like `go` routines, `makeChan`, `select`, `defer`, `ref`/`deref`, and `make`/`cap`.
- **Smarter Type Linker**: Automatic, recursive type generation for Go structs, interfaces, and methods. Supports struct embedding (`extends`) and nested type resolution.
- **True Parallelism**: Goroutine-based workers with zero-copy shared memory (`typego:memory`).
- **Modern Package Ecosystem**: Built-in CLI for managing Go dependencies with `typego.modules.json` and `typego.lock`.
- **Fast Developer Loop**: Hot-reloading dev server and ~0.2s interpreter startup. Compiles to single-binary with `--compile`.

---

## Quick Start

### Installation

```bash
go install github.com/repyh/typego/cmd/typego@latest
```

### Usage

**1. Initialize a Project**

```bash
typego init myapp
cd myapp
```

**2. Run Development Server**

```bash
# Watch mode with hot-reload
typego dev src/index.ts
```

**3. Build for Production**

```bash
# Compile to a standalone binary
typego build src/index.ts -o app
./app
```

---

## Language Reference

### Imports

TypeGo uses a unique import scheme to distinguish between TypeScript/JS modules and Go packages.

```typescript
// Import standard Go packages
// Access any package from the Go standard library by prefixing with 'go:'
import { Println, Sprintf, Errorf } from "go:fmt";
import { Sleep, Mutex } from "go:sync";

// Import external Go modules (must be added via CLI first)
// Use the full module path prefixed with 'go:'
import { Default } from "go:github.com/gin-gonic/gin";

// Import TypeGo internal modules
// Access runtime-specific features like workers and shared memory
import { Worker } from "typego:worker";
import { makeShared } from "typego:memory";

// Import relative TypeScript files
// Standard ES module imports work as expected
import { util } from "./util";
import config from "./config.json";
```

### Intrinsics

TypeGo exposes low-level Go primitives as global functions, allowing you to write "Go-like" TypeScript. These are always available in the global scope.

#### Core Primitives

| Function | Signature | Description |
|----------|-----------|-------------|
| `go` | `go(fn: () => void)` | Launches a background goroutine. The function runs concurrently. |
| `makeChan` | `makeChan<T>(size?: number)` | Creates a synchronized Go channel. Optional buffer size. |
| `select` | `select(cases: SelectCase[])` | Multiplexes channel operations. Blocks until one case proceeds. |
| `defer` | `defer(fn: () => void)` | Registers a function to run when the current `typego.scope` exits. |
| `panic` | `panic(err: any)` | Triggers a native Go panic. Can be caught by Go's `recover` or JS `try/catch`. |
| `recover` | `recover()` | Recovers from a panic inside a `defer` block. Returns the error or `null`. |
| `iota` | `const` | Auto-incrementing compile-time constant (useful for enums). |

#### Memory & Slices

| Function | Signature | Description |
|----------|-----------|-------------|
| `ref` | `ref<T>(val: T): Ref<T>` | Creates a pointer handle to a value on the Go heap. |
| `deref` | `deref<T>(ptr: Ref<T>): T` | Dereferences a pointer or `Ref` object to get its value. |
| `make` | `make(Type, len, cap?)` | Allocates high-performance slices (TypedArrays). Maps to Go's `make`. |
| `cap` | `cap(v: any): number` | Returns the capacity of a slice, channel, or buffer. |
| `copy` | `copy(dst, src): number` | Performs high-speed memory copying between buffers/slices. |
| `sizeof` | `sizeof(obj): number` | Estimates the memory footprint of a JS/Go object in bytes. |

#### Process & Environment

Access system-level process information globally via the `process` object.

```typescript
// Environment Variables
console.log(process.env.HOME);
console.log(process.env.NODE_ENV); // Emulated

// Current Working Directory
console.log(process.cwd());

// Platform Information
console.log(process.platform); // e.g., 'linux', 'darwin', 'windows'
console.log(process.version);  // Go runtime version (e.g., 'go1.21.0')

// Arguments
console.log(process.argv);     // Command line arguments
```

#### Encoding

Built-in support for fast string/byte conversion, compatible with the Web `TextEncoder`/`TextDecoder` API.

```typescript
const encoder = new TextEncoder();
const decoder = new TextDecoder();

const str = "Hello TypeGo";
const bytes = encoder.encode(str); // Returns Uint8Array
const decoded = decoder.decode(bytes); // Returns string

// Direct hex/base64 utilities (if available) would be documented here
```

#### IO Utilities

Helpers for bridging JS objects with Go's `io.Reader` and `io.Writer` interfaces.

```typescript
// Wrap a JS object with a read() method into an io.Reader
const reader = wrapReader({
  read: (buf) => {
    // Fill buf with data, return number of bytes read
    return 0;
  }
});

// Wrap a JS object with a write() method into an io.Writer
const writer = wrapWriter({
  write: (buf) => {
    // Process buf, return number of bytes written
    return buf.byteLength;
  }
});
```

---

### Standard Library

TypeGo includes pre-bound versions of common Go standard library packages. These are highly optimized for performance.

#### go:fmt

Print formatted output to stdout/stderr.

```typescript
import { Println, Printf, Sprintf, Errorf } from "go:fmt";

// Basic printing
Println("Hello", "World");

// Formatted printing
Printf("User %s has ID %d\n", "Alice", 123);

// String formatting
const msg = Sprintf("Value: %.2f", 3.14159);

// Creating errors
const err = Errorf("failed to connect: %s", "timeout");
```

#### go:net/http

Make HTTP requests or run a server. Supports the standard Go `http.Client` and `http.Server`.

```typescript
import { Get, Post, ListenAndServe, Fetch } from "go:net/http";

// --- Client ---

// Simple GET
const resp = Get("https://api.example.com/data");
console.log(`Status: ${resp.Status}`);
console.log(`Body: ${resp.Body}`); // Body is automatically read as string

// POST with JSON
const res = Post("https://api.example.com/users", JSON.stringify({ name: "Bob" }));
console.log(res.StatusCode);

// Fetch API (Promise-based wrapper)
Fetch("https://example.com").then(r => {
  console.log("Async fetch complete:", r.Status);
});

// --- Server ---

// Start a simple HTTP server
const server = ListenAndServe(":8080", (w, r) => {
    // w: ResponseWriter, r: Request
    // Note: This callback signature adapts Go's handler interface
    console.log("Received request:", r.Method, r.URL.Path);

    // Write response (conceptual)
    // w.Write(encoder.encode("Hello from TypeGo!"));
});

// Graceful shutdown
// server.close();
```

#### go:os

File system interactions (sandboxed to CWD by default for security).

```typescript
import * as os from "go:os";

// Reading and Writing files
try {
    os.WriteFile("config.txt", "port=8080");
    const content = os.ReadFile("config.txt");
    console.log("Config:", content);
} catch (e) {
    console.error("File error:", e);
}

// Directory management
os.Mkdir("logs");
os.MkdirAll("data/users/images");

// Cleanup
os.Remove("config.txt");
os.RemoveAll("data");

// Environment and Process
const path = Getenv("PATH");
if (Args.length > 1) {
    console.log("First arg:", Args[1]);
}

// Exit the program
// Exit(1);
```

#### go:sync

Concurrency primitives.

```typescript
import { Spawn, Sleep, Chan } from "go:sync";

// --- Spawning Async Tasks ---
Spawn(async () => {
    console.log("Starting background task...");
    await Sleep(1000); // Sleep for 1000ms
    console.log("Task complete");
});

// --- Channels ---
const ch = new Chan();
ch.send("Hello");
const msg = await ch.recv();
console.log(msg);
```

#### typego:memory

Shared memory for workers. This module provides mechanisms to share raw memory between isolated runtime instances (workers).

```typescript
import { makeShared } from "typego:memory";

// Create a named shared buffer (1KB)
// This buffer is allocated in Go memory and mapped to both runtimes
const sharedBuf = makeShared("globalCounter", 1024);

const view = new Int32Array(sharedBuf.buffer);
Atomics.add(view, 0, 1); // Thread-safe atomic operation
```

#### typego:worker

Thread-based worker spawning. Each worker runs in its own goroutine with a separate Sobek runtime, ensuring true parallelism.

```typescript
import { Worker } from "typego:worker";

// Spawn a new worker from a file
const w = new Worker("./worker.ts");

// Send data to the worker
w.postMessage({ command: "process", data: [1, 2, 3] });

// Receive messages from the worker
w.onmessage = (msg) => {
    console.log("Worker replied:", msg);
};

// Terminate the worker
// w.terminate();
```

---

### Concurrency

TypeGo offers true parallelism via goroutines, distinct from Node.js's single-threaded event loop.

#### Goroutines

Use the `go` keyword to launch a lightweight thread.

```typescript
import { Println } from "go:fmt";

go(() => {
    // This runs concurrently!
    // Blocking operations here do NOT block the main thread
    Println("Background work");
});
```

#### Channels

Channels provide a safe way for goroutines to communicate.

```typescript
// Create a buffered channel of numbers
const ch = makeChan<number>(5);

go(() => {
    // Send values into the channel
    ch.send(1);
    ch.send(2);
    ch.close(); // Close the channel when done
});

// Receive values
// Note: In TypeGo, receive() is blocking if not in a select/async context
// Use select for non-blocking or multiplexed operations
```

#### Select

The `select` statement lets a goroutine wait on multiple communication operations.

```typescript
import { Println } from "go:fmt";

const ch1 = makeChan<string>();
const ch2 = makeChan<string>();

go(() => { ch1.send("one"); });
go(() => { ch2.send("two"); });

// Select blocks until one case is ready
select([
  {
    chan: ch1,
    recv: (v) => Println("Received from ch1:", v)
  },
  {
    chan: ch2,
    recv: (v) => Println("Received from ch2:", v)
  },
  {
    default: () => Println("No messages ready")
  }
]);
```

---

### Memory Management

TypeGo manages memory automatically but provides tools for manual control when performance is critical.

#### Pointers & Refs

When passing large Go structs to JavaScript, TypeGo typically copies the value. To avoid this overhead, use `ref`.

```typescript
// 'ref' and 'deref' are global intrinsics

const largeObj = { id: 1, payload: "huge..." };
const ptr = ref(largeObj); // ptr is a light handle (Ref<T>)

// Pass 'ptr' around cheaply...

const val = deref(ptr); // Access the value when needed
```

#### Shared Memory

Use `typego:memory` to share buffers between workers without serialization.

```typescript
import { makeShared } from "typego:memory";

// Main thread
const shared = makeShared("pixels", 1920 * 1080 * 4); // { buffer, mutex } with 8MB buffer

// Create a view over the shared buffer
const pixels = new Uint8ClampedArray(shared.buffer);
// When mutating shared state, coordinate using shared.mutex

// Worker thread
// Receive the same { buffer, mutex } (e.g. via message from main or a global binding)
// const workerPixels = new Uint8ClampedArray(shared.buffer);
// Use shared.mutex here as well when writing
```

#### Defer

Ensure resources are cleaned up when a scope exits.

```typescript
import { WriteFile, Remove } from "go:os";

function processFile() {
    const path = "data.tmp";
    WriteFile(path, "temporary data");

    defer(() => {
        Remove(path);
        console.log("Temp file removed");
    });

    // Process file...
    // If a panic occurs here, Remove() is still called
}
```

---

## Tooling

### CLI Reference

| Command | Arguments | Description |
|---------|-----------|-------------|
| `run` | `<file>` | Executes a TypeScript file using the fast interpreter mode. Does not produce a binary. |
| `dev` | `<file>` | Starts a development server that watches for file changes and hot-reloads the application. |
| `build` | `<file> [-o output]` | Compiles the TypeScript entrypoint and all dependencies into a standalone executable. |
| `init` | `[name]` | Scaffolds a new TypeGo project. Creates `typego.modules.json`, `package.json`, and directory structure. |
| `types` | None | Generates `.d.ts` definition files for all configured Go imports. |
| `add` | `<pkg>` | Adds a Go module dependency to `typego.modules.json` and runs resolution. |
| `remove` | `<pkg>` | Removes a Go module dependency. |
| `list` | None | Lists all configured Go dependencies and their versions. |
| `update` | `[pkg]` | Updates all or specific Go modules to their latest versions. |
| `outdated` | None | Checks for newer versions of configured Go modules. |
| `install` | None | Manually triggers the JIT build and dependency resolution process. |
| `clean` | None | Cleans the `.typego/` workspace, removing cached artifacts and types. |

### Package Management

TypeGo uses `typego.modules.json` to manage Go dependencies.

**Example `typego.modules.json`:**

```json
{
  "modules": {
    "github.com/gin-gonic/gin": "v1.9.1",
    "github.com/go-redis/redis/v8": "v8.11.5"
  }
}
```

**Adding a dependency:**

```bash
typego add github.com/gin-gonic/gin
```

This command:
1. Updates `typego.modules.json`.
2. Runs `go get`.
3. Rebuilds the internal bridge.
4. Generates TypeScript definitions (`go.d.ts`).

---

## Ecosystem

### Performance

TypeGo is optimized for high-throughput I/O and true parallelism.

- **JIT Compilation**: TypeScript is compiled to JS via esbuild, then executed by Sobek (Goja).
- **Zero-Copy**: Intrinsics use direct memory access where possible.
- **Goroutines**: `go()` spawns lightweight threads, allowing thousands of concurrent tasks.

For a detailed breakdown of the execution model, interop overhead, and optimization strategies, see **[OPTIMIZATION.md](OPTIMIZATION.md)**.

### Runtime Comparison

| Feature | TypeGo | Node.js | Deno | Bun |
|---------|:------:|:-------:|:----:|:---:|
| **Language** | Go + TS | C++ + JS | Rust + JS | Zig + JS |
| **Concurrency** | Goroutines (M:N) | Event Loop (1:1) | Tokio (M:N)* | Event Loop |
| **Parallelism** | True (Multi-core) | Single-threaded | Worker-based | Worker-based |
| **Binary Size** | ~10MB | ~80MB | ~100MB | ~90MB |
| **Startup** | Fast (<0.2s) | Medium | Fast | Instant |
| **Ecosystem** | Go Modules | NPM | Deno/NPM | NPM |

*\*Deno uses Tokio but JS runs on V8 isolates which are effectively single-threaded per isolate.*

---

## Cookbook

### HTTP Server with Middleware

```typescript
import { ListenAndServe } from "go:net/http";
import { Println } from "go:fmt";

function loggingMiddleware(next) {
    return (w, r) => {
        const start = Date.now();
        next(w, r);
        Println(r.Method, r.URL.Path, `${Date.now() - start}ms`);
    };
}

const handler = (w, r) => {
    if (r.URL.Path === "/") {
        // w.Write("Welcome!");
        Println("Root visited");
    } else {
        // w.WriteHeader(404);
    }
};

const srv = ListenAndServe(":3000", loggingMiddleware(handler));
```

### Worker Pool

Implement a pool of workers processing tasks from a shared channel.

```typescript
import { Worker } from "typego:worker";
import { Println } from "go:fmt";

// 'makeChan' and 'go' are global intrinsics

const jobs = makeChan<number>(100);
const results = makeChan<number>(100);

// Worker logic (inlined for demo, usually in separate file)
function startWorker(id: number) {
    const w = new Worker("./worker_script.ts");
    w.postMessage({ id, jobsChan: jobs, resultsChan: results });
    // Note: Passing channels to workers requires specialized serialization support
    // For this example, assume workers pull from a shared source or use shared memory
}

// Start 3 workers
for (let i = 1; i <= 3; i++) {
    startWorker(i);
}

// Send jobs
for (let j = 1; j <= 5; j++) {
    jobs.send(j);
}
jobs.close();

// Collect results
for (let a = 1; a <= 5; a++) {
   // const res = results.recv();
   // Println("Result:", res);
}
```

### File Processing Pipeline

Read a file, transform it, and write it back using buffers.

```typescript
import { ReadFile, WriteFile } from "go:os";

const data = ReadFile("input.dat");
// TODO: transform `data` here if needed

WriteFile("output.dat", data);
console.log(`Copied ${data.length} bytes`);
```

### External Go Modules

Using a Redis client (`go-redis`).

```typescript
// 1. typego add github.com/go-redis/redis/v8
import { NewClient } from "go:github.com/go-redis/redis/v8";
import { context } from "go:context";

const ctx = context.Background();
const rdb = NewClient({
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
});

const err = rdb.Set(ctx, "key", "value", 0).Err();
if (err) panic(err);

const val = rdb.Get(ctx, "key").Val();
console.log("key:", val);
```

---

## Troubleshooting

**Problem: "Module not found"**
- Ensure you have run `typego add <module>`.
- Run `typego install` to force a JIT rebuild.
- Check `typego.modules.json`.

**Problem: "Undefined symbol in Go module"**
- Some complex Go types (generics, complex structs) may not map 1:1 to JS.
- Check `.typego/go.d.ts` to see what was generated.
- Use `any` casting as a workaround if types are too strict.

**Problem: Performance is slow**
- Avoid frequent calls across the Go-JS boundary in tight loops.
- Use `TypedArrays` and `ArrayBuffers` instead of JS arrays for large data.
- Use `go:sync` intrinsics instead of busy-waiting in JS.

---

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
./typego run examples/02-concurrency-basics.ts
./typego run examples/09-typego-stdlib.ts
```

---

## License

MIT
