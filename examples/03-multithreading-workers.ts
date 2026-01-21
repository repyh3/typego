import { Spawn } from "go:sync";
import { Println } from "go:fmt";
import { makeShared } from "go:memory";

/**
 * Multi-Threading & Shared Memory Showcase
 * 
 * Demonstrates:
 * 1. Worker threads for parallel execution.
 * 2. Shared buffers for zero-copy state sharing.
 */

// 1. Create a shared buffer (1 KB)
const shared = makeShared("main_state", 1024);
(shared as any)[0] = 0; // Initialize a counter at index 0

async function startWorker(id: number) {

    // In a real app, you would use SpawnWorker(jsFile)
    Spawn(async () => {
        for (let i = 0; i < 10; i++) {
            (shared as any)[0]++;
            Println(`Worker ${id} incremented counter to: ${(shared as any)[0]}`);
        }
    });
}

Println("🏭 Spawning workers...");
startWorker(1);
startWorker(2);

Println("Main thread yields...");
