// MODULE: typego:worker
declare module "typego:worker" {
    export interface Worker {
        postMessage(msg: any): void;
        terminate(): void;
        onmessage: (msg: { data: any }) => void;
    }
    export var Worker: {
        new(scriptPath: string): Worker;
    }
}
// END: typego:worker
