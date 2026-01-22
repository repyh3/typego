import { Println } from "go:fmt";

/**
 * Memory Management and Slices
 * 
 * Demonstrates:
 * 1. 'sizeof' - Deep memory inspection of JS objects.
 * 2. 'make' - Go-style allocation with length and capacity.
 * 3. 'cap' - Inspecting allocation capacity.
 * 4. 'copy' - High-performance memory copying.
 */

function main() {
    Println("=== Memory Management Showcase ===");

    // 1. Size Inspection
    const obj = { a: 1, b: "hello", c: new Uint8Array(100) };
    Println(`Size of object: ${sizeof(obj)} bytes`);

    // 2. Allocation with Capacity
    // make(Type, length, capacity)
    const buffer = make(Uint8Array, 10, 100);
    Println(`Buffer length: ${buffer.length}, capacity: ${cap(buffer)}`);

    // 3. Efficient Copying
    const src = new Uint8Array([1, 2, 3, 4, 5]);
    const dst = new Uint8Array(3);

    // copy returns number of elements copied
    const n = copy(dst, src);
    Println(`Copied ${n} elements. Dest: [${Array.from(dst).join(", ")}]`);

    // 4. Working with buffers
    const large = make(Uint8Array, 0, 1024);
    Println(`Empty buffer with 1KB capacity: cap=${cap(large)}`);
}

main();
