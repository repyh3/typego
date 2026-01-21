// Package polyfills provides Node.js-compatible globals for the Goja runtime
package polyfills

import (
	"github.com/dop251/goja"
	"github.com/repyh3/typego/eventloop"
)

func EnableAll(vm *goja.Runtime, el *eventloop.EventLoop) {
	EnableProcess(vm)
	EnableBuffer(vm)
	EnableTimers(vm, el)
	EnableEncoding(vm)
}
