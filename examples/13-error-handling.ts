import { Println } from "go:fmt";

/**
 * Advanced Control Flow: Defer, Panic, and Recover
 * 
 * Demonstrates:
 * 1. 'defer' - Scheduling cleanup tasks.
 * 2. 'panic' - Triggering a fatal error (like throw but with Go semantics).
 * 3. 'recover' - Capturing panics safely.
 */

function safeDivision(a: number, b: number): number {
    defer(() => Println("  [Cleanup] Division operation concluded."));

    if (b === 0) {
        panic("Division by zero!");
    }
    return a / b;
}

function main() {
    Println("=== Error Handling Showcase ===");

    defer(() => {
        const err = recover();
        if (err) {
            Println(`🔥 Successfully recovered from panic: ${err}`);
        }
    });

    Println("Attempting 10 / 2...");
    const res = safeDivision(10, 2);
    Println(`Result: ${res}`);

    Println("\nAttempting 10 / 0...");
    safeDivision(10, 0);

    Println("This line will not be reached.");
}

main();
