/**
 * TextEncoder takes a stream of code points as input and emits a stream of UTF-8 bytes.
 */
interface TextEncoder {
    readonly encoding: string;
    /**
     * Encodes the given string into a Uint8Array.
     */
    encode(input?: string): Uint8Array;
}

declare var TextEncoder: {
    prototype: TextEncoder;
    new(): TextEncoder;
};

/**
 * TextDecoder represents a decoder for a specific text encoding, such as UTF-8.
 */
interface TextDecoder {
    /**
     * The encoding used by this decoder.
     */
    readonly encoding: string;
    /**
     * Decodes the given input buffer into a string.
     */
    decode(input?: ArrayBufferView | ArrayBuffer): string;
}

declare var TextDecoder: {
    prototype: TextDecoder;
    new(label?: string, options?: { fatal?: boolean; ignoreBOM?: boolean }): TextDecoder;
};

