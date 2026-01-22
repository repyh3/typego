/**
 * Launches a new goroutine to execute the provided function.
 * @param fn The function to execute in the background.
 * @param args Arguments to pass to the function.
 */
declare function go<T extends (...args: any[]) => any>(fn: T, ...args: Parameters<T>): void;

/**
 * A Go-style channel for synchronized communication.
 */
interface Chan<T> {
    /**
     * Sends a value into the channel. Blocks if the channel is full.
     */
    send(val: T): void;

    /**
     * Receives a value from the channel. Blocks if the channel is empty.
     */
    recv(): T;

    /**
     * Closes the channel. Subsequent sends will panic.
     */
    close(): void;
}

/**
 * Creates a new channel with an optional buffer size.
 * @param size The buffer size (0 for unbuffered).
 */
declare function makeChan<T>(size?: number): Chan<T>;

/**
 * Multiplexes multiple channel operations.
 * Implements Go's select statement logic.
 */
declare function select(cases: Array<{
    chan?: Chan<any>;
    send?: any;
    recv?: (val: any) => void;
    case?: () => void;
    default?: () => void;
}>): void;
