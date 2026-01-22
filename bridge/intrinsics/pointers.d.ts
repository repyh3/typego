/**
 * Creates a reference (pointer) to a value on the Go heap.
 * This allows passing primitives by reference and manual memory management.
 * 
 * @param val The value to reference.
 * @returns A Ref object with 'value' (get/set) and 'ptr' (uintptr).
 */
declare function ref<T>(val: T): { value: T; readonly ptr: number };

/**
 * Dereferences a raw pointer or a Ref object.
 * 
 * @param ptr The raw uintptr or Ref object.
 * @returns The value stored at that address.
 */
declare function deref<T>(ptr: number | { readonly ptr: number }): T;
