import { Println } from "go:fmt";

/**
 * Low-Level Pointers
 * 
 * Demonstrates:
 * 1. 'ref' - Getting a reference to a variable.
 * 2. 'deref' - Resolving a reference to its current value.
 */

function main() {
    Println("=== Pointer Intrinsics Showcase ===");

    let secret = "Initial Value";

    // Create a pointer to 'secret'
    const ptr = ref(secret);

    Println(`Current value via deref: ${deref(ptr)}`);

    // Modify the original variable
    secret = "Updated Value";

    // The pointer tracks the change because it points to the variable's "slot"
    Println(`New value via same pointer: ${deref(ptr)}`);

    // Pointers work through scopes as long as the variable is alive
    const update = (p: any) => {
        secret = "Modified in closure";
    };

    update(ptr);
    Println(`After closure update: ${deref(ptr)}`);
}

main();
