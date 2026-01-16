import { WriteFile, ReadFile } from "go/os";
import { Println } from "go/fmt";

/**
 * Secure Sandbox (The Vault) Showcase
 * 
 * Demonstrates Phase 2 Hardening:
 * 1. Path Jailing: Files can only be accessed within the project root.
 * 2. Symlink Protection: Evaluates real paths to prevent escaping the jail.
 */

async function main() {
    Println("🔒 Testing Secure Sandbox...");

    // 1. Valid access within the jail
    Println("\n--- Valid Access ---");
    try {
        WriteFile("sandbox_demo.txt", "TypeGo Protection: ACTIVE");
        const data = ReadFile("sandbox_demo.txt");
        Println(`✅ Read successful: "${data}"`);
    } catch (e) {
        Println(`❌ Unexpected failure: ${e}`);
    }

    // 2. Blocked access outside the jail (Relative escape)
    Println("\n--- Blocked Escape (../../) ---");
    try {
        WriteFile("../../jailbreak.txt", "escaping...");
        Println("❌ FAILURE: Escaped the jail!");
    } catch (e) {
        Println(`✅ Success: Blocked. Error: ${e}`);
    }

    // 3. Blocked access outside the jail (Absolute escape)
    Println("\n--- Blocked Escape (C:/Windows/win.ini) ---");
    try {
        ReadFile("C:/Windows/win.ini");
        Println("❌ FAILURE: Read outside jail!");
    } catch (e) {
        Println(`✅ Success: Blocked. Error: ${e}`);
    }
}

main();
