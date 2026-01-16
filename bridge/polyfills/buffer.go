package polyfills

import "github.com/dop251/goja"

// BufferShimJS is the JavaScript polyfill for the Buffer global
const BufferShimJS = `if (typeof Buffer === 'undefined') { 
	globalThis.Buffer = { 
		from: function(str) { 
			if (typeof str === 'string') { 
				var arr = []; 
				for (var i = 0; i < str.length; i++) { 
					arr.push(str.charCodeAt(i)); 
				} 
				return new Uint8Array(arr); 
			} 
			return new Uint8Array(str); 
		}, 
		alloc: function(size) { return new Uint8Array(size); },
		toString: function() { return '[object Buffer]'; }
	}; 
}`

// EnableBuffer injects the Buffer global via JavaScript
func EnableBuffer(vm *goja.Runtime) {
	_, _ = vm.RunString(BufferShimJS)
}
