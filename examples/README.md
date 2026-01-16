# TypeGo Showcase 🚀

This directory contains examples that demonstrate the core capabilities and security features of the TypeGo runtime.

## Run Examples

You can run any of these examples using the TypeGo CLI:

```bash
typego run examples/<filename>.ts
```

For examples involving memory limits:
```bash
typego run examples/04-secure-network.ts -M 256
```

---

## Showcase Scripts

### 1. `01-hello-world.ts`
The basic "Hello World" demonstrating TypeScript compilation, type safety, and standard output bridging.

### 2. `02-go-concurrency.ts`
Demonstrates how TypeGo leverages Go's concurrency. Includes `Spawn` for background tasks and asynchronous `Sleep`.

### 3. `03-multithreading-workers.ts`
Showcases high-performance multi-threading using the **Worker API** and **SharedArrayBuffer** equivalents for zero-copy state sharing between Go and JS.

### 4. `04-secure-network.ts`
Demonstrates the hardened HTTP bridge. Showcases:
- 30s timeouts.
- 50MB response size limits.
- Promise-based `Fetch` API.

### 5. `05-secure-fs.ts`
Demonstrates the **Secure Sandbox (The Vault)**. Showcases how File System access is jailed to the workspace root and protected against symlink-based escapes.

### 6. `06-process-env.ts`
Showcases the security-filtered environment variables. Demonstrates how sensitive information (like `AWS_SECRET`) is automatically stripped, while whitelisted variables (like `PATH`) are preserved.

---

## Advanced: The Cluster Demo
The `examples/cluster` directory contains a full-scale demonstration of:
- A hybrid Go/TypeScript HTTP server.
- Worker pools for intensive computation.
- Shared mutexes and state across workers.