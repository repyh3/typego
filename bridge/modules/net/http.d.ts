// MODULE: go:net/http
declare module "go:net/http" {
    export function ListenAndServe(addr: string, handler: (req: Request, res: Response) => Promise<void> | void): void;
    export function Get(url: string): { Status: string; StatusCode: number; Body: string };
    export function Post(url: string, body: string, contentType?: string): Promise<{ Status: string; StatusCode: number; Body: string }>;
    export function Fetch(url: string): Promise<{ Status: string; StatusCode: number; Body: string }>;

    export interface Request {
        Method: string;
        URL: {
            Path: string;
            RawQuery: string;
            Host: string;
        };
        Host: string;
        Proto: string;
        Header: Record<string, string[]>;
        Body: any;
        ContentLength: number;
        RemoteAddr: string;
    }

    export interface Response {
        setHeader(key: string, value: string): void;
        status(code: number): Response;
        write(data: string): Response;
        send(data?: string): void;
        json(data: any): void;
        redirect(url: string, code?: number): void;
    }

    export interface Server {
        close(timeout?: number): void;
    }
}
// END: go:net/http
