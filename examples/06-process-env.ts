import { Println } from "go/fmt";

/**
 * Environment Isolation Showcase
 * 
 * Demonstrates secure environment filtering in Phase 2:
 * 1. Sensitive vars (APPDATA, AWS_SECRET, etc.) are stripped.
 * 2. Whitelisted system vars (PATH, LANG) are preserved.
 * 3. Custom 'TYPEGO_' prefixed vars are allowed.
 */

async function main() {
    Println("🧠 Checking process.env isolation...");

    // System whitelisted
    Println(`PATH exists: ${process.env.PATH ? "✅" : "❌"}`);

    // Potentially sensitive (Automatically blocked by our whitelist-only policy)
    const sensitiveKeys = ["APPDATA", "USERPROFILE", "AWS_SECRET", "DB_PASSWORD"];
    for (const key of sensitiveKeys) {
        const leaked = process.env[key];
        Println(`${key}: ${leaked ? "❌ (LEAKED!)" : "✅ (Blocked)"}`);
    }

    // Custom TypeGo variables
    Println(`\nNote: Any environment variable prefixed with 'TYPEGO_' is automatically allowed.`);
    Println(`Current TYPEGO_ variables: ${Object.keys(process.env).filter(k => k.startsWith("TYPEGO_")).join(", ") || "None"}`);
}

main();
