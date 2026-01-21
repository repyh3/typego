import { Default, Context } from "go:github.com/gin-gonic/gin";

// === In-memory link storage ===
const links: Map<string, { url: string; clicks: number; createdAt: Date }> = new Map();

function saveLink(code: string, url: string): void {
    links.set(code, { url, clicks: 0, createdAt: new Date() });
}

function getLink(code: string): { url: string; clicks: number } | null {
    const link = links.get(code);
    return link ? { url: link.url, clicks: link.clicks } : null;
}

function incrementClicks(code: string): void {
    const link = links.get(code);
    if (link) link.clicks++;
}

function getAllLinks(): Array<{ code: string; url: string; clicks: number }> {
    return Array.from(links.entries()).map(([code, data]) => ({
        code,
        url: data.url,
        clicks: data.clicks,
    }));
}

// === Short ID generator ===
const ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
function generateId(): string {
    let result = "";
    for (let i = 0; i < 7; i++) {
        result += ALPHABET[Math.floor(Math.random() * ALPHABET.length)];
    }
    return result;
}

// === Handlers ===
function handleShorten(c: Context): void {
    // For now, use query params since body parsing needs investigation
    const url = c.Query("url");

    if (!url) {
        c.JSON(400, { error: "Missing 'url' query parameter. Use ?url=https://example.com" });
        return;
    }

    if (!url.startsWith("http://") && !url.startsWith("https://")) {
        c.JSON(400, { error: "URL must start with http:// or https://" });
        return;
    }

    const code = generateId();
    saveLink(code, url);

    c.JSON(201, {
        code,
        short_url: `http://localhost:8080/${code}`,
        original_url: url,
    });
}

function handleRedirect(c: Context): void {
    const code = c.Param("code");

    // Skip API routes
    if (code === "api") {
        c.Next();
        return;
    }

    const link = getLink(code);
    if (!link) {
        c.JSON(404, { error: "Short link not found", code });
        return;
    }

    incrementClicks(code);
    c.Redirect(302, link.url);
}

function handleStats(c: Context): void {
    const code = c.Param("code");
    const link = getLink(code);

    if (!link) {
        c.JSON(404, { error: "Short link not found" });
        return;
    }

    c.JSON(200, {
        code,
        original_url: link.url,
        clicks: link.clicks,
    });
}

// === Server setup ===
const app = Default();

// Health check
app.GET("/api/health", (c: Context) => {
    c.JSON(200, { status: "ok", service: "shortlink", links_count: links.size });
});

// Link management
app.POST("/api/shorten", handleShorten);
app.GET("/api/shorten", handleShorten); // Also allow GET for easy testing
app.GET("/api/stats/:code", handleStats);
app.GET("/api/links", (c: Context) => {
    c.JSON(200, { links: getAllLinks() });
});

// Redirect (must be last)
app.GET("/:code", handleRedirect);

console.log("🔗 Shortlink API v1.0");
console.log("   Endpoints:");
console.log("   GET  /api/shorten?url=...  - Create short link");
console.log("   GET  /:code                - Redirect to original");
console.log("   GET  /api/stats/:code      - View click stats");
console.log("   GET  /api/links            - List all links");
console.log("");
console.log("🚀 Starting on http://localhost:8080");

app.Run(":8080");
