/**
 * Panics the runtime with a message, mimicking Go's panic.
 * It stops the ordinary flow of control and begins panicking.
 * 
 * @param message The panic message.
 */
declare function panic(message: string): never;
