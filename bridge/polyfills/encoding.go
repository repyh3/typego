package polyfills

import "github.com/dop251/goja"

// EncodingShimJS provides TextEncoder and TextDecoder polyfills for the Goja runtime.
// These are Web APIs commonly used for string/byte conversion.
const EncodingShimJS = `
(function() {
	if (typeof TextEncoder === 'undefined') {
		globalThis.TextEncoder = function TextEncoder() {};
		TextEncoder.prototype.encode = function(str) {
			if (typeof str !== 'string') {
				str = String(str);
			}
			var arr = [];
			for (var i = 0; i < str.length; i++) {
				var code = str.charCodeAt(i);
				if (code < 0x80) {
					arr.push(code);
				} else if (code < 0x800) {
					arr.push(0xC0 | (code >> 6));
					arr.push(0x80 | (code & 0x3F));
				} else if (code < 0x10000) {
					arr.push(0xE0 | (code >> 12));
					arr.push(0x80 | ((code >> 6) & 0x3F));
					arr.push(0x80 | (code & 0x3F));
				}
			}
			return new Uint8Array(arr);
		};
	}

	if (typeof TextDecoder === 'undefined') {
		globalThis.TextDecoder = function TextDecoder(encoding) {
			this.encoding = encoding || 'utf-8';
		};
		TextDecoder.prototype.decode = function(input) {
			if (!input) return '';
			var bytes = input instanceof Uint8Array ? input : new Uint8Array(input);
			var result = '';
			var i = 0;
			while (i < bytes.length) {
				var byte1 = bytes[i++];
				if (byte1 < 0x80) {
					result += String.fromCharCode(byte1);
				} else if ((byte1 & 0xE0) === 0xC0) {
					var byte2 = bytes[i++] & 0x3F;
					result += String.fromCharCode(((byte1 & 0x1F) << 6) | byte2);
				} else if ((byte1 & 0xF0) === 0xE0) {
					var byte2 = bytes[i++] & 0x3F;
					var byte3 = bytes[i++] & 0x3F;
					result += String.fromCharCode(((byte1 & 0x0F) << 12) | (byte2 << 6) | byte3);
				}
			}
			return result;
		};
	}
})();
`

// EnableEncoding injects TextEncoder and TextDecoder globals.
func EnableEncoding(vm *goja.Runtime) {
	_, _ = vm.RunString(EncodingShimJS)
}
