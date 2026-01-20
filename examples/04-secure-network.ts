import { Get } from "go:net/http";
import { Println } from "go:fmt";

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
        // Go's http.Get is synchronous, but we can treat it as such or await if wrapped.
        // TypeGo generated functions are synchronous in JS unless they return Promise.
        // Since http.Get blocks, it will return the Response object directly.
        const res = Get("https://httpbin.org/get?msg=TypeGo_Secure_Fetch");
        Println(`Status: ${res.Status}`);

    } catch (e) {
        Println(`❌ Fetch Error: ${e}`);
    }

    // Attempting a larger request (This will fail due to guardrails)
    Println("\n🛡️ Test: Fetching oversized data...");
    try {
        const res = Get("https://httpbin.org/bytes/60000000");
        Println(`Status: ${res.Status}`);
    } catch (e) {
        Println(`✅ Success: Error caught: ${e}`);
    }
}

testFetch();
