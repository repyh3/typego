export function main() {
    console.log("Testing native defer syntax...");

    defer(() => console.log("Deferred log 1 (last)"));
    defer(() => console.log("Deferred log 2 (first)"));

    console.log("Running main body");
}

main();