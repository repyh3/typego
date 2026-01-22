import { Sleep } from "go:sync";
import { Println } from "go:fmt";

/**
 * Concurrency Basics
 * 
 * Demonstrates:
 * 1. 'go' - The native Go goroutine intrinsic.
 * 2. 'Sleep' - Non-blocking sleep that yields to the event loop.
 */

async function main() {
    Println("🚀 Starting concurrency basics...");

    // Launch a background goroutine
    go(async () => {
        for (let i = 0; i < 5; i++) {
            await Sleep(500);
            Println(`  [Goroutine] Pulse ${i + 1}`);
        }
    });

    Println("Main: Doing work...");
    await Sleep(1000);
    Println("Main: Still working...");
    await Sleep(1000);

    Println("🏁 Main thread finished.");
}

main();
