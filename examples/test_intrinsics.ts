export function main() {
    console.log("Testing intrinsics...");

    // Test sizeof
    console.log("sizeof(bool):", sizeof(true)); // Expect 1
    console.log("sizeof(int):", sizeof(42));    // Expect 8

    // Test panic
    try {
        panic("simulated panic");
    } catch (e: any) {
        console.log("Caught panic:", e.message);
    }

    // Test scope/defer
    console.log("Testing scope...");
    typego.scope((defer) => {
        defer(() => console.log("defer 1 (last)"));
        defer(() => console.log("defer 2 (first)"));
        console.log("Inside scope");
    });
    console.log("After scope");
}

main();