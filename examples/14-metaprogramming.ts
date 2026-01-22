import { Println } from "go:fmt";

/**
 * Metaprogramming with iota
 * 
 * Demonstrates:
 * 1. 'iota' - Auto-incrementing compile-time constants.
 */

// iota starts at 0 and increments per declaration in the same file.
const READ = iota;     // 0
const WRITE = iota;    // 1
const EXECUTE = iota;  // 2

const STATUS_IDLE = iota;    // 3
const STATUS_RUNNING = iota; // 4

function main() {
    Println("=== iota Metaprogramming Showcase ===");

    Println("Permissions:");
    Println(`  READ:    ${READ}`);
    Println(`  WRITE:   ${WRITE}`);
    Println(`  EXECUTE: ${EXECUTE}`);

    Println("\nSystem Status:");
    Println(`  IDLE:    ${STATUS_IDLE}`);
    Println(`  RUNNING: ${STATUS_RUNNING}`);
}

main();
