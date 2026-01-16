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

	/**
	 * ReadFile reads the named file and returns the contents.
	 * In TypeGo, this is sandboxed to the project root.
	 */
	export function ReadFile(filename: string): string;
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




// MODULE: go:color
declare module "go:github.com/fatih/color" {
	/**
	 * New returns a newly created color object.
	 */
	export function New(...args: any[]): any;
	/**
	 * RGB returns a new foreground color in 24-bit RGB.
	 */
	export function RGB(...args: any[]): any;
	/**
	 * BgRGB returns a new background color in 24-bit RGB.
	 */
	export function BgRGB(...args: any[]): any;
	/**
	 * Set sets the given parameters immediately. It will change the color of
	 * output with the given SGR parameters until color.Unset() is called.
	 */
	export function Set(...args: any[]): any;
	/**
	 * Unset resets all escape attributes and clears the output. Usually should
	 * be called after Set().
	 */
	export function Unset(...args: any[]): any;
	/**
	 * Black is a convenient helper function to print with black foreground. A
	 * newline is appended to format by default.
	 */
	export function Black(...args: any[]): any;
	/**
	 * Red is a convenient helper function to print with red foreground. A
	 * newline is appended to format by default.
	 */
	export function Red(...args: any[]): any;
	/**
	 * Green is a convenient helper function to print with green foreground. A
	 * newline is appended to format by default.
	 */
	export function Green(...args: any[]): any;
	/**
	 * Yellow is a convenient helper function to print with yellow foreground.
	 * A newline is appended to format by default.
	 */
	export function Yellow(...args: any[]): any;
	/**
	 * Blue is a convenient helper function to print with blue foreground. A
	 * newline is appended to format by default.
	 */
	export function Blue(...args: any[]): any;
	/**
	 * Magenta is a convenient helper function to print with magenta foreground.
	 * A newline is appended to format by default.
	 */
	export function Magenta(...args: any[]): any;
	/**
	 * Cyan is a convenient helper function to print with cyan foreground. A
	 * newline is appended to format by default.
	 */
	export function Cyan(...args: any[]): any;
	/**
	 * White is a convenient helper function to print with white foreground. A
	 * newline is appended to format by default.
	 */
	export function White(...args: any[]): any;
	/**
	 * BlackString is a convenient helper function to return a string with black
	 * foreground.
	 */
	export function BlackString(...args: any[]): any;
	/**
	 * RedString is a convenient helper function to return a string with red
	 * foreground.
	 */
	export function RedString(...args: any[]): any;
	/**
	 * GreenString is a convenient helper function to return a string with green
	 * foreground.
	 */
	export function GreenString(...args: any[]): any;
	/**
	 * YellowString is a convenient helper function to return a string with yellow
	 * foreground.
	 */
	export function YellowString(...args: any[]): any;
	/**
	 * BlueString is a convenient helper function to return a string with blue
	 * foreground.
	 */
	export function BlueString(...args: any[]): any;
	/**
	 * MagentaString is a convenient helper function to return a string with magenta
	 * foreground.
	 */
	export function MagentaString(...args: any[]): any;
	/**
	 * CyanString is a convenient helper function to return a string with cyan
	 * foreground.
	 */
	export function CyanString(...args: any[]): any;
	/**
	 * WhiteString is a convenient helper function to return a string with white
	 * foreground.
	 */
	export function WhiteString(...args: any[]): any;
	/**
	 * HiBlack is a convenient helper function to print with hi-intensity black foreground. A
	 * newline is appended to format by default.
	 */
	export function HiBlack(...args: any[]): any;
	/**
	 * HiRed is a convenient helper function to print with hi-intensity red foreground. A
	 * newline is appended to format by default.
	 */
	export function HiRed(...args: any[]): any;
	/**
	 * HiGreen is a convenient helper function to print with hi-intensity green foreground. A
	 * newline is appended to format by default.
	 */
	export function HiGreen(...args: any[]): any;
	/**
	 * HiYellow is a convenient helper function to print with hi-intensity yellow foreground.
	 * A newline is appended to format by default.
	 */
	export function HiYellow(...args: any[]): any;
	/**
	 * HiBlue is a convenient helper function to print with hi-intensity blue foreground. A
	 * newline is appended to format by default.
	 */
	export function HiBlue(...args: any[]): any;
	/**
	 * HiMagenta is a convenient helper function to print with hi-intensity magenta foreground.
	 * A newline is appended to format by default.
	 */
	export function HiMagenta(...args: any[]): any;
	/**
	 * HiCyan is a convenient helper function to print with hi-intensity cyan foreground. A
	 * newline is appended to format by default.
	 */
	export function HiCyan(...args: any[]): any;
	/**
	 * HiWhite is a convenient helper function to print with hi-intensity white foreground. A
	 * newline is appended to format by default.
	 */
	export function HiWhite(...args: any[]): any;
	/**
	 * HiBlackString is a convenient helper function to return a string with hi-intensity black
	 * foreground.
	 */
	export function HiBlackString(...args: any[]): any;
	/**
	 * HiRedString is a convenient helper function to return a string with hi-intensity red
	 * foreground.
	 */
	export function HiRedString(...args: any[]): any;
	/**
	 * HiGreenString is a convenient helper function to return a string with hi-intensity green
	 * foreground.
	 */
	export function HiGreenString(...args: any[]): any;
	/**
	 * HiYellowString is a convenient helper function to return a string with hi-intensity yellow
	 * foreground.
	 */
	export function HiYellowString(...args: any[]): any;
	/**
	 * HiBlueString is a convenient helper function to return a string with hi-intensity blue
	 * foreground.
	 */
	export function HiBlueString(...args: any[]): any;
	/**
	 * HiMagentaString is a convenient helper function to return a string with hi-intensity magenta
	 * foreground.
	 */
	export function HiMagentaString(...args: any[]): any;
	/**
	 * HiCyanString is a convenient helper function to return a string with hi-intensity cyan
	 * foreground.
	 */
	export function HiCyanString(...args: any[]): any;
	/**
	 * HiWhiteString is a convenient helper function to return a string with hi-intensity white
	 * foreground.
	 */
	export function HiWhiteString(...args: any[]): any;
}
// END: go:color

// MODULE: go:fmt
declare module "go:fmt" {
	/**
	 * Errorf formats according to a format specifier and returns the string as a
	 * value that satisfies error.
	 * 
	 * If the format specifier includes a %w verb with an error operand,
	 * the returned error will implement an Unwrap method returning the operand.
	 * If there is more than one %w verb, the returned error will implement an
	 * Unwrap method returning a []error containing all the %w operands in the
	 * order they appear in the arguments.
	 * It is invalid to supply the %w verb with an operand that does not implement
	 * the error interface. The %w verb is otherwise a synonym for %v.
	 */
	export function Errorf(...args: any[]): any;
	/**
	 * FormatString returns a string representing the fully qualified formatting
	 * directive captured by the [State], followed by the argument verb. ([State] does not
	 * itself contain the verb.) The result has a leading percent sign followed by any
	 * flags, the width, and the precision. Missing flags, width, and precision are
	 * omitted. This function allows a [Formatter] to reconstruct the original
	 * directive triggering the call to Format.
	 */
	export function FormatString(...args: any[]): any;
	/**
	 * Fprintf formats according to a format specifier and writes to w.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Fprintf(...args: any[]): any;
	/**
	 * Printf formats according to a format specifier and writes to standard output.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Printf(...args: any[]): any;
	/**
	 * Sprintf formats according to a format specifier and returns the resulting string.
	 */
	export function Sprintf(...args: any[]): any;
	/**
	 * Appendf formats according to a format specifier, appends the result to the byte
	 * slice, and returns the updated slice.
	 */
	export function Appendf(...args: any[]): any;
	/**
	 * Fprint formats using the default formats for its operands and writes to w.
	 * Spaces are added between operands when neither is a string.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Fprint(...args: any[]): any;
	/**
	 * Print formats using the default formats for its operands and writes to standard output.
	 * Spaces are added between operands when neither is a string.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Print(...args: any[]): any;
	/**
	 * Sprint formats using the default formats for its operands and returns the resulting string.
	 * Spaces are added between operands when neither is a string.
	 */
	export function Sprint(...args: any[]): any;
	/**
	 * Append formats using the default formats for its operands, appends the result to
	 * the byte slice, and returns the updated slice.
	 */
	export function Append(...args: any[]): any;
	/**
	 * Fprintln formats using the default formats for its operands and writes to w.
	 * Spaces are always added between operands and a newline is appended.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Fprintln(...args: any[]): any;
	/**
	 * Println formats using the default formats for its operands and writes to standard output.
	 * Spaces are always added between operands and a newline is appended.
	 * It returns the number of bytes written and any write error encountered.
	 */
	export function Println(...args: any[]): any;
	/**
	 * Sprintln formats using the default formats for its operands and returns the resulting string.
	 * Spaces are always added between operands and a newline is appended.
	 */
	export function Sprintln(...args: any[]): any;
	/**
	 * Appendln formats using the default formats for its operands, appends the result
	 * to the byte slice, and returns the updated slice. Spaces are always added
	 * between operands and a newline is appended.
	 */
	export function Appendln(...args: any[]): any;
	/**
	 * Scan scans text read from standard input, storing successive
	 * space-separated values into successive arguments. Newlines count
	 * as space. It returns the number of items successfully scanned.
	 * If that is less than the number of arguments, err will report why.
	 */
	export function Scan(...args: any[]): any;
	/**
	 * Scanln is similar to [Scan], but stops scanning at a newline and
	 * after the final item there must be a newline or EOF.
	 */
	export function Scanln(...args: any[]): any;
	/**
	 * Scanf scans text read from standard input, storing successive
	 * space-separated values into successive arguments as determined by
	 * the format. It returns the number of items successfully scanned.
	 * If that is less than the number of arguments, err will report why.
	 * Newlines in the input must match newlines in the format.
	 * The one exception: the verb %c always scans the next rune in the
	 * input, even if it is a space (or tab etc.) or newline.
	 */
	export function Scanf(...args: any[]): any;
	/**
	 * Sscan scans the argument string, storing successive space-separated
	 * values into successive arguments. Newlines count as space. It
	 * returns the number of items successfully scanned. If that is less
	 * than the number of arguments, err will report why.
	 */
	export function Sscan(...args: any[]): any;
	/**
	 * Sscanln is similar to [Sscan], but stops scanning at a newline and
	 * after the final item there must be a newline or EOF.
	 */
	export function Sscanln(...args: any[]): any;
	/**
	 * Sscanf scans the argument string, storing successive space-separated
	 * values into successive arguments as determined by the format. It
	 * returns the number of items successfully parsed.
	 * Newlines in the input must match newlines in the format.
	 */
	export function Sscanf(...args: any[]): any;
	/**
	 * Fscan scans text read from r, storing successive space-separated
	 * values into successive arguments. Newlines count as space. It
	 * returns the number of items successfully scanned. If that is less
	 * than the number of arguments, err will report why.
	 */
	export function Fscan(...args: any[]): any;
	/**
	 * Fscanln is similar to [Fscan], but stops scanning at a newline and
	 * after the final item there must be a newline or EOF.
	 */
	export function Fscanln(...args: any[]): any;
	/**
	 * Fscanf scans text read from r, storing successive space-separated
	 * values into successive arguments as determined by the format. It
	 * returns the number of items successfully parsed.
	 * Newlines in the input must match newlines in the format.
	 */
	export function Fscanf(...args: any[]): any;
}
// END: go:fmt
