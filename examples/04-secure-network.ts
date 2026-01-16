import { Fetch } from "go/net/http";
import { Println } from "go/fmt";

/**
 * Secure Network Showcase
 * 
 * Demonstrates Phase 1 Hardening:
 * 1. 30s connection timeouts.
 * 2. 50MB response body limit (prevents OOM attacks).
 * 3. Promise-based Fetch API.
 */

async function testFetch() {
    Println("🌐 Fetching from API...");

    try {
        const res = await Fetch("https://httpbin.org/get?msg=TypeGo_Secure_Fetch");
        Println(`Status: ${res.StatusCode} ${res.Status}`);

        const bodySnippet = res.Body.substring(0, 100);
        Println(`Response Body (snippet): ${bodySnippet}...`);

    } catch (e) {
        Println(`❌ Fetch Error: ${e}`);
    }

    // Attempting a larger request (This will fail due to guardrails)
    Println("\n🛡️ Testing Guardrails: Fetching oversized data (Limit: 50MB)...");
    try {
        // Requesting 60MB (This should trigger the 50MB guardrail)
        await Fetch("https://httpbin.org/bytes/60000000");
        Println("❌ Failed: Guardrail didn't catch large response!");
    } catch (e) {
        Println(`✅ Success: Guardrail blocked large response. Error: ${e}`);
    }
}

testFetch();
