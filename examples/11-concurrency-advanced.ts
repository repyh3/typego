import { Sleep } from "go:sync";
import { Println } from "go:fmt";

/**
 * Advanced Concurrency: Channels and Select
 * 
 * Demonstrates:
 * 1. 'makeChan' - Creating Go-backed communication channels.
 * 2. 'select' - Multiplexing channel operations (send/recv).
 * 3. 'go' - Parallel execution with shared memory.
 */

async function main() {
    Println("=== Advanced Concurrency Showcase ===");

    const c1 = makeChan<string>();
    const c2 = makeChan<string>();

    // Producer 1
    go(async () => {
        await Sleep(1000);
        c1.send("Message from Alpha");
    });

    // Producer 2
    go(async () => {
        await Sleep(500);
        c2.send("Message from Beta");
    });

    // We expect to receive from C2 first because it sleeps for less time
    for (let i = 0; i < 2; i++) {
        select([
            {
                chan: c1,
                recv: (val) => Println(`[RECEIVE] Channel 1: ${val}`)
            },
            {
                chan: c2,
                recv: (val) => Println(`[RECEIVE] Channel 2: ${val}`)
            },
            {
                default: () => {
                    // This runs if No channel is ready
                    // In this loop, we might see this if we didn't await or if logic was different
                }
            }
        ]);
    }

    Println("🏁 Done.");
}

main();
