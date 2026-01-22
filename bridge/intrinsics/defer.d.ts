/**
 * Schedules a function call to be run immediately before the function returns.
 * The deferred call's arguments are evaluated immediately, but the function call 
 * is not executed until the surrounding function returns.
 * 
 * Note: This is a TypeGo intrinsic handled by the AST transformer.
 * 
 * @param fn The function to execute when the surrounding function exits.
 */
declare function defer(fn: () => void): void;
