package intrinsics

// BufferShimJS is the JavaScript polyfill for the Buffer global.
const BufferShimJS = `if (typeof Buffer === 'undefined') { 
	globalThis.Buffer = { 
		from: function(data, encoding) { 
			return __bufferFrom(data, encoding);
		}, 
		alloc: function(size) { 
			return __bufferAlloc(size);
		},
		toString: function() { return '[object Buffer]'; }
	}; 
}`

// EncodingShimJS provides TextEncoder and TextDecoder polyfills.
const EncodingShimJS = `
(function() {
	if (typeof TextEncoder === 'undefined') {
		globalThis.TextEncoder = function TextEncoder() {};
		TextEncoder.prototype.encode = function(str) {
			return __encode(str);
		};
	}

	if (typeof TextDecoder === 'undefined') {
		globalThis.TextDecoder = function TextDecoder(encoding) {
			this.encoding = encoding || 'utf-8';
		};
		TextDecoder.prototype.decode = function(input) {
			return __decode(input);
		};
	}
})();
`

// EnableGlobals injects all environment globals into the VM.
// This replaces the legacy polyfills package.
func (r *Registry) EnableGlobals() {
	// 1. Native backing intrinsics (already set in Enable)

	// 2. JS Shims for Standard APIs
	_, _ = r.vm.RunString(EncodingShimJS)
	_, _ = r.vm.RunString(BufferShimJS)

	// 3. Environment Globals
	r.EnableProcess()
	r.EnableTimers()
}
