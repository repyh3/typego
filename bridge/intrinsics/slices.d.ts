/**
 * Returns the capacity of a typed array, array, or buffer.
 * For JS Arrays, this is the same as length.
 * For TypedArrays, this is the allocated size of the underlying buffer.
 */
declare function cap(v: any[] | ArrayBuffer | SharedArrayBuffer | { length: number } | { byteLength: number }): number;

/**
 * Creates a new slice (TypedArray) with a specified length and optional capacity.
 * @param Type The TypedArray constructor (e.g. Uint8Array).
 * @param len The initial length.
 * @param cap The allocated capacity (must be >= len).
 */
declare function make<T extends { new(len: number): any }>(Type: T, len: number, cap?: number): InstanceType<T>;

/**
 * Copies elements from a source to a destination.
 * Returns the number of elements copied.
 */
declare function copy(dst: { set(src: any): void; length: number }, src: { length: number }): number;
