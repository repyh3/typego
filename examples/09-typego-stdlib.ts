import { makeShared, stats, ptr } from "typego:memory";
import { Worker } from "typego:worker";
import { Println, Printf } from "go:fmt";
import { Sleep } from "go:sync";

Println("--- TypeGo Stdlib Verification ---");

// 1. Test Memory Module
Println("1. Testing typego:memory...");
const s = stats();
Printf("Current Alloc: %d\n", s.alloc);

const shared = makeShared("test-buffer", 1024);
const u8 = new Uint8Array(shared.buffer);
u8[0] = 42;
Printf("Shared Buffer[0] initialized to: %d\n", u8[0]);

// 2. Test Worker Module
Println("\n2. Testing typego:worker...");

try {
    const w = new Worker("examples/01-hello-world.ts");
    Println("✅ Worker created successfully.");
    w.terminate();
} catch (e) {
    Printf("❌ Worker creation failed: %v\n", e);
}

Println("\n✅ Verification Complete.");
