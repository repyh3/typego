package intrinsics

import (
	"io"
	"github.com/grafana/sobek"
)

type JSReader struct {
	vm   *sobek.Runtime
	obj  *sobek.Object
	read sobek.Callable
}

func (r *JSReader) Read(p []byte) (n int, err error) {
	buf := r.vm.NewArrayBuffer(p)
	res, err := r.read(r.obj, r.vm.ToValue(buf))
	if err != nil {
		return 0, err
	}
	
	n = int(res.ToInteger())
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (r *Registry) WrapReader(call sobek.FunctionCall) sobek.Value {
	obj := call.Argument(0).ToObject(r.vm)
	if obj == nil {
		return sobek.Undefined()
	}
	
	read, ok := sobek.AssertFunction(obj.Get("read"))
	if !ok {
		panic(r.vm.NewTypeError("object must have a read(buffer) method"))
	}
	
	return r.vm.ToValue(&JSReader{vm: r.vm, obj: obj, read: read})
}

type JSWriter struct {
	vm    *sobek.Runtime
	obj   *sobek.Object
	write sobek.Callable
}

func (w *JSWriter) Write(p []byte) (n int, err error) {
	buf := w.vm.NewArrayBuffer(p)
	res, err := w.write(w.obj, w.vm.ToValue(buf))
	if err != nil {
		return 0, err
	}
	return int(res.ToInteger()), nil
}

func (r *Registry) WrapWriter(call sobek.FunctionCall) sobek.Value {
	obj := call.Argument(0).ToObject(r.vm)
	if obj == nil {
		return sobek.Undefined()
	}

	write, ok := sobek.AssertFunction(obj.Get("write"))
	if !ok {
		panic(r.vm.NewTypeError("object must have a write(buffer) method"))
	}

	return r.vm.ToValue(&JSWriter{vm: r.vm, obj: obj, write: write})
}
