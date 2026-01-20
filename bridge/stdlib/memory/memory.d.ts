// MODULE: go:memory
declare module "go:memory" {
    export function makeShared(name: string, size: number): Uint8Array;
}
// END: go:memory

// MODULE: typego:memory
declare module "typego:memory" {
    export function makeShared(name: string, size: number): { buffer: ArrayBuffer; mutex: any };
    export function stats(): { alloc: number; totalAlloc: number; sys: number; numGC: number };
    export function ptr(val: any): any;
}
// END: typego:memory
