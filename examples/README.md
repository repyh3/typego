# TypeGo Examples

A collection of examples showcasing TypeGo's unique capabilities, from Go-style concurrency to low-level memory management.

## Getting Started

To run any example:
```bash
typego run examples/<filename>.ts
```

## Available Examples

### Basics
- `01-hello-world.ts`: The classic entry point.
- `02-concurrency-basics.ts`: Simple background tasks using `go` and `Sleep`.
- `06-process-env.ts`: Accessing environment variables.
- `09-typego-stdlib.ts`: Overview of built-in modules.

### Concurrency & Parallelism
- `11-concurrency-advanced.ts`: Channels (`makeChan`) and multiplexing (`select`).
- `03-multithreading-workers.ts`: High-performance background workers.

### Security & Sandbox
- `04-secure-network.ts`: Granular network permissions.
- `05-secure-fs.ts`: Strict file system access control.

### Low-Level & Memory
- `12-memory-management.ts`: Detailed inspection with `sizeof` and allocation with `make`/`cap`.
- `15-pointers.ts`: Variable references and tracking using `ref` and `deref`.
- `08-nested-structs.ts`: Complex data structure bridging.

### Advanced Control & Meta
- `13-error-handling.ts`: Robust cleanup with `defer` and panic recovery.
- `14-metaprogramming.ts`: Auto-incrementing constants with `iota`.
- `07-external-module.ts`: Importing and using external TypeScript modules.

### Networking
- `10-http-server.ts`: Building a type-safe web server with the `net/http` module.