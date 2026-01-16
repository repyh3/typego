
import { Println } from "go:fmt";

try {
    Println("Accessing process.env:", process.env["OS"] || process.env["USER"] || "Unknown");
    Println("Process Platform:", process.platform);
} catch (e: any) {
    Println("Process error:", e.message);
}

try {
    const b = Buffer.from("Hello Buffer");
    Println("Buffer content:", b.toString());
} catch (e: any) {
    Println("Buffer error:", e.message);
}
