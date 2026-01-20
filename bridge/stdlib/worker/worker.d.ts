// MODULE: typego:worker
declare module "typego:worker" {
    export class Worker {
        constructor(scriptPath: string);
        postMessage(msg: any): void;
        terminate(): void;
        onmessage: (msg: { data: any }) => void;
    };
}
// END: typego:worker
