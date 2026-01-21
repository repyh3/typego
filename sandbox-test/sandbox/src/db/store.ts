// In-memory link storage
const links: Map<string, { url: string; clicks: number; createdAt: Date }> = new Map();

export function saveLink(code: string, url: string): void {
    links.set(code, { url, clicks: 0, createdAt: new Date() });
}

export function getLink(code: string): { url: string; clicks: number } | null {
    const link = links.get(code);
    return link ? { url: link.url, clicks: link.clicks } : null;
}

export function incrementClicks(code: string): void {
    const link = links.get(code);
    if (link) link.clicks++;
}

export function getAllLinks(): Array<{ code: string; url: string; clicks: number }> {
    return Array.from(links.entries()).map(([code, data]) => ({
        code,
        url: data.url,
        clicks: data.clicks,
    }));
}
