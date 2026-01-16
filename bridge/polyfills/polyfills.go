// Package polyfills provides Node.js-compatible globals for the Goja runtime
package polyfills

import (
	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

// EnableAll injects all Node.js polyfills into the VM
func EnableAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	EnableProcess(vm)
	EnableBuffer(vm)
	EnableTimers(vm, el)
}
