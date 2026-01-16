import { Println } from "go/fmt";
import { Sleep } from "go/sync";

async function main() {
    Println("🛡️ Starting Memory Guardrail Test...");

    try {
        const sink: any[] = [];
        Println("Allocating memory rapidly...");

        // This loop should be interrupted by the engine's memory monitor
        while (true) {
            sink.push(new Array(1000000).fill("A".repeat(1024))); // Massive allocation
            if (sink.length % 1 === 0) {
                Println(`Current sink size: ${sink.length}`);
            }
        }
    } catch (e) {
        Println(`EXPECTED CATCH: ${e}`);
    }
}

main();
