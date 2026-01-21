import { ListenAndServe, Request, Response } from "go:net/http";
import { Println } from "go:fmt";

const PORT = ":8080";

Println(`Starting server on http://localhost${PORT}`);

ListenAndServe(PORT, async (req: Request, res: Response) => {
    Println(`[${req.method}] ${req.url}`);

    // CORS
    res.setHeader("Access-Control-Allow-Origin", "*");
    res.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE");
    res.setHeader("Content-Type", "application/json");

    if (req.method === "OPTIONS") {
        res.status(204).send();
        return;
    }

    try {
        if (req.path === "/") {
            res.json({
                message: "Welcome to TypeGo HTTP Server",
                endpoints: ["/echo", "/time", "/error"]
            });
            return;
        }

        if (req.path === "/time") {
            res.json({
                time: new Date().toISOString(),
                timestamp: Date.now()
            });
            return;
        }

        if (req.path === "/echo") {
            if (req.method === "POST") {
                const body = await req.body();
                let data = body;
                try {
                    data = JSON.parse(body);
                } catch (e) {
                    // keep as string
                }

                res.json({
                    received: data,
                    headers: req.headers,
                    query: req.query
                });
            } else {
                res.status(405).json({ error: "Method not allowed. Use POST." });
            }
            return;
        }

        if (req.path === "/error") {
            throw new Error("Something went wrong!");
        }

        res.status(404).json({ error: "Not Found", path: req.path });

    } catch (e) {
        if (isGoError(e)) {
            Println(`Handler error: ${e.message}`);
            res.status(500).json({ error: e.message });
        } else {
            Println(`Unknown error: ${e}`);
            res.status(500).json({ error: "Internal Server Error" });
        }
    }
});
