import { Println } from "go/fmt";

/**
 * TypeGo Hello World
 * 
 * Demonstrates:
 * 1. TypeScript compilation and execution.
 * 2. Standard Go output bridging via 'go/fmt'.
 */

const message: string = "Hello from TypeGo TypeScript!";
Println(message);

const stats = {
    version: "v1.0.0",
    engine: "Goja (Go-based JS Engine)",
    type: "Safe/Isolated"
};

Println(`Runtime: ${stats.engine} - ${stats.version}`);
