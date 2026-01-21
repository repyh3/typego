# TypeGo Examples

Demonstration scripts showcasing TypeGo's core features.

## Usage

```bash
typego run examples/<filename>.ts
```

With memory limit:
```bash
typego run examples/04-secure-network.ts -M 256
```

## Examples

| File | Description |
|------|-------------|
| `01-hello-world.ts` | Basic TypeScript execution and output |
| `02-go-concurrency.ts` | `Spawn` and async `Sleep` for parallelism |
| `03-multithreading-workers.ts` | Worker API with shared memory |
| `04-secure-network.ts` | HTTP fetch with timeouts and size limits |
| `05-secure-fs.ts` | Sandboxed file system access |
| `06-process-env.ts` | Filtered environment variables |
| `07-external-module.ts` | Third-party Go package imports |
| `08-nested-structs.ts` | Complex nested struct handling |
| `09-typego-stdlib.ts` | Native `typego:memory` and `typego:worker` |
| `10-http-server.ts` | HTTP server with CORS and routing |

## NPM Dependencies

Some examples use NPM packages. Install with:

```bash
cd examples
npm install
```