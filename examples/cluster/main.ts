import chalk from "chalk";
import { Println } from "go:fmt";

Println(chalk.blue.bold("🚀 TypeGo Cluster Demo"));

let processedCount = 0;

// Simple UUID alternative (no crypto dependency)
function simpleId(): string {
    return Math.random().toString(36).substring(2, 10);
}

async function processRequests() {
    for (let i = 0; i < 5; i++) {
        const requestId = simpleId();

        Println(chalk.yellow(`📨 [${requestId}] Processing Task #${i + 1}...`));

        let hash = 0;
        for (let j = 0; j < 100000; j++) {
            hash = (hash * 31 + j) % 1000000007;
        }

        processedCount++;

        Println(chalk.green(`✅ [${requestId}] Complete -> Hash: ${hash}`));

        await new Promise((r) => setTimeout(r, 50));
    }

    Println(chalk.magenta.bold(`\n📊 Total Requests Processed: ${processedCount}`));
    Println(chalk.cyan("🏁 Cluster Demo Complete!"));
}

processRequests().catch((e) => Println("Error: " + e));
