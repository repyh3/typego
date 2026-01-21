// Package polyfills provides Node.js-compatible globals for the Goja runtime
package polyfills

import (
	"github.com/dop251/goja"
	"github.com/repyh/typego/eventloop"
)

// EnableAll injects all polyfills into the VM:
//   - process: Environment variables, platform info
//   - Buffer: from(), alloc() for buffer operations
//   - Timers: setTimeout, setInterval, clearInterval
//   - Encoding: TextEncoder, TextDecoder for string/byte conversion
func EnableAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	EnableProcess(vm)
	EnableBuffer(vm)
	EnableTimers(vm, el)
	EnableEncoding(vm)
}
