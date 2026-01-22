/**
 * Go-style auto-incrementing constant.
 * Only works during compile-time inside 'const' declarations.
 * Resets to 0 at the start of each file.
 */
declare const iota: number;
