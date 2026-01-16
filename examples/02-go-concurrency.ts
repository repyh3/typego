import { Spawn, Sleep } from "go/sync";
import { Println } from "go/fmt";

/**
 * Go Concurrency Showcase
 * 
 * Demonstrates:
 * 1. 'Spawn' - Running background tasks in Go goroutines from JS.
 * 2. 'Sleep' - Asynchronous, non-blocking sleep that yields to the event loop.
 */

async function main() {
    Println("🚀 Starting concurrency demo...");

    // Spawn a background "heartbeat" task
    Spawn(async () => {
        for (let i = 0; i < 5; i++) {
            await Sleep(500);
            Println(`  [Heartbeat] Pulse ${i + 1}`);
        }
    });

    // Run a sequence
    Println("Main: Doing work...");
    await Sleep(1000);
    Println("Main: Still working...");
    await Sleep(1000);

    Println("🏁 Main thread finished.");
}

main();
