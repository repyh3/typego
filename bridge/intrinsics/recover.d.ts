/**
 * Recovers from a panic, mimicking Go's recover.
 * Can only be called inside a deferred function.
 * 
 * @returns The value passed to panic(), or null if no panic is active.
 */
declare function recover(): any;
