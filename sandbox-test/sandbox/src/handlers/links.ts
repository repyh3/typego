import { Context } from "go:github.com/gin-gonic/gin";
import { saveLink, getLink } from "../db/store";
import { generate } from "../utils/shortid";

interface ShortenRequest {
    url: string;
}

export function handleShorten(c: Context): void {
    const body = c.Request.Body;
    // Parse JSON body
    let req: ShortenRequest;
    try {
        req = JSON.parse(body);
    } catch {
        c.JSON(400, { error: "Invalid JSON body" });
        return;
    }

    if (!req.url) {
        c.JSON(400, { error: "Missing 'url' field" });
        return;
    }

    // Basic URL validation
    if (!req.url.startsWith("http://") && !req.url.startsWith("https://")) {
        c.JSON(400, { error: "URL must start with http:// or https://" });
        return;
    }

    const code = generate();
    saveLink(code, req.url);

    const host = c.Request.Host || "localhost:8080";
    c.JSON(201, {
        code,
        short_url: `http://${host}/${code}`,
        original_url: req.url,
    });
}

export function handleRedirect(c: Context): void {
    const code = c.Param("code");

    // Skip API routes
    if (code === "api") {
        c.Next();
        return;
    }

    const link = getLink(code);
    if (!link) {
        c.JSON(404, { error: "Short link not found" });
        return;
    }

    // Async click increment (fire and forget)
    import("../db/store").then(store => store.incrementClicks(code));

    c.Redirect(302, link.url);
}

export function handleStats(c: Context): void {
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
