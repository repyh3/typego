/**
 * Estimates the memory size of an object in bytes.
 * This is a TypeGo intrinsic that performs shallow or deep inspection 
 * based on the underlying Go representation.
 * 
 * @param obj The object to measure.
 * @returns Estimated size in bytes.
 */
declare function sizeof(obj: any): number;
