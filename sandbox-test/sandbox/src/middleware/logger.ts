import { Context } from "go:github.com/gin-gonic/gin";

export function requestLogger(c: Context): void {
    const start = Date.now();
    const method = c.Request.Method;
    const path = c.Request.URL.Path;

    c.Next();

    const duration = Date.now() - start;
    const status = c.Writer.Status();
    console.log(`[${method}] ${path} - ${status} (${duration}ms)`);
}
