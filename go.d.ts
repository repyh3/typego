declare module "go/memory" {
    /**
     * Ptr simulates a Go pointer by wrapping a value in an object.
     */
    export class Ptr<T> {
        constructor(value: T);
        value: T;
    }

    export interface SharedBuffer {
        buffer: Uint8Array;
        mutex: import("go/sync").Mutex;
    }

    /**
     * makeShared creates a named shared memory segment.
     */
    export function makeShared(name: string, size: number): SharedBuffer;
}

declare module "go/sync" {
    /**
     * Mutex provides async Read-Write locking
     */
    export interface Mutex {
        lock(): Promise<void>;
        unlock(): void;
        rlock(): Promise<void>;
        runlock(): void;
    }
    /**
     * Spawn schedules a function to run. 
     * In this runtime, it ensures the function runs on the JavaScript thread.
     */
    export function Spawn(fn: () => void): void;

    /**
     * Sleep pauses the current execution for the specified milliseconds.
     * Returns a promise that resolves when the time is up.
     */
    export function Sleep(ms: number): Promise<void>;

    /**
     * Chan is a Go-style channel for communication between tasks.
     */
    export class Chan<T> {
        constructor();
        send(val: T): void;
        recv(): Promise<T>;
    }
}

declare module "go/fmt" {
    /**
     * Println formats using the default formats for its operands and writes to standard output.
     */
    export function Println(...args: any[]): void;
}

declare module "go/os" {
    /**
     * WriteFile writes data to the named file.
     * In TypeGo, this is sandboxed to the project root.
     */
    export function WriteFile(filename: string, data: string): void;
}

declare module "go/net/http" {
    export interface Response {
        Status: string;
        StatusCode: number;
        Body: string;
    }

    /**
     * Get issues a GET to the specified URL.
     */
    export function Get(url: string): Response;

    /**
     * Fetch issues an async GET to the specified URL.
     */
    export function Fetch(url: string): Promise<Response>;
}

/**
 * CLI-provided supercharged tools
 */
declare const native: {
    StartTime: string;
    GetRuntimeInfo(): string;
};

/**
 * A 1KB shared memory buffer provided by the CLI
 */
declare const cliBuffer: Uint8Array;

/**
 * Access Go-level memory statistics.
 */
declare function goMemory(): {
    alloc: number;
    totalAlloc: number;
    sys: number;
    numGC: number;
};

declare class Worker {
    constructor(scriptPath: string);
    postMessage(msg: any): void;
    terminate(): void;
    onmessage: ((event: { data: any }) => void) | null;
}

declare const self: {
    postMessage(msg: any): void;
    onmessage: ((event: { data: any }) => void) | null;
};



