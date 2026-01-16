import { Fetch } from "go/net/http";
import { Println } from "go/fmt";

async function main() {
    Println("🛡️ Starting HTTP Guardrail Test...");

    try {
        Println("Testing Response Size Limit (fetching large file)...");
        // This should hit the 50MB limit
        // Using a known large file or a service that generates large payload
        const resp = await Fetch("https://httpbin.org/bytes/100000000"); // 100MB
        Println(`Success? Status: ${resp.Status}`);
    } catch (e) {
        Println(`EXPECTED ERROR (Size): ${e}`);
    }

    try {
        Println("\nTesting Timeout (fetching slow url)...");
        // This should hit the 30s timeout
        const resp = await Fetch("https://httpbin.org/delay/40"); // 40s delay
        Println(`Success? Status: ${resp.Status}`);
    } catch (e) {
        Println(`EXPECTED ERROR (Timeout): ${e}`);
    }
}

main();
