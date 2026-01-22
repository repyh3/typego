/**
 * Buffer provides a way of handling binary data in TypeGo, mirrored from Node.js.
 * In TypeGo, Buffer is a subclass of Uint8Array with additional performance optimizations.
 */
declare interface Buffer extends Uint8Array { }

declare var Buffer: {
    /**
     * Allocates a new Buffer of size bytes.
     * Backed by high-performance Go memory allocation.
     */
    alloc(size: number): Buffer;
    /**
     * Creates a new Buffer from the given data.
     * If data is a string, it defaults to UTF-8 encoding.
     */
    from(data: string | number[] | Uint8Array | ArrayBuffer, encoding?: string): Buffer;
};
