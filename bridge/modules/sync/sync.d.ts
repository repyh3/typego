// MODULE: go:sync
declare module "go:sync" {
    export function Spawn(fn: () => void): void;
    export function Sleep(ms: number): Promise<void>;
    export function Chan<T = any>(buffer?: number): {
        Send(val: T): void;
        Recv(): T;
        TryRecv(): [T, boolean];
        Close(): void;
    };
}
// END: go:sync
