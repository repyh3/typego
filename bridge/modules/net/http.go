// Package http provides bindings for Go's net/http package.
package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grafana/sobek"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
)

func init() {
	core.RegisterModule(&httpModule{})
}

type httpModule struct {
	el *eventloop.EventLoop
}

func (m *httpModule) Name() string {
	return "go:net/http"
}

func (m *httpModule) Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	m.el = el
	Register(vm, el)
}

// Default HTTP client with production-ready timeouts
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

const maxResponseBodySize = 50 * 1024 * 1024 // 50MB

type Module struct{}

func (h *Module) Get(vm *sobek.Runtime) func(sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		url := call.Argument(0).String()

		resp, err := httpClient.Get(url)
		if err != nil {
			panic(vm.NewTypeError(fmt.Sprintf("http.Get error: %v", err)))
		}
		defer resp.Body.Close()

		limit := int64(maxResponseBodySize)
		body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
		if err != nil && err != io.EOF {
			panic(vm.NewTypeError(fmt.Sprintf("http body read error: %v", err)))
		}

		if int64(len(body)) > limit {
			panic(vm.NewTypeError(fmt.Sprintf("http response too large (max %d MB)", maxResponseBodySize/1024/1024)))
		}

		res := vm.NewObject()
		_ = res.Set("Status", resp.Status)
		_ = res.Set("StatusCode", resp.StatusCode)
		_ = res.Set("Body", string(body))

		return res
	}
}

func Register(vm *sobek.Runtime, el *eventloop.EventLoop) {
	h := &Module{}
	server := NewServer(vm, el)

	obj := vm.NewObject()
	_ = obj.Set("Get", h.Get(vm))

	_ = obj.Set("Post", func(call sobek.FunctionCall) sobek.Value {
		url := call.Argument(0).String()
		body := call.Argument(1).String()
		contentType := "application/json"
		if len(call.Arguments) > 2 {
			contentType = call.Argument(2).String()
		}

		p, resolve, reject := el.CreatePromise()

		go func() {
			req, err := http.NewRequest("POST", url, io.NopCloser(
				io.LimitReader(
					&stringReader{s: body, i: 0},
					int64(len(body)),
				),
			))
			if err != nil {
				el.RunOnLoop(func() {
					reject(vm.NewGoError(err))
				})
				return
			}
			req.Header.Set("Content-Type", contentType)

			resp, err := httpClient.Do(req)
			if err != nil {
				el.RunOnLoop(func() {
					reject(vm.NewGoError(err))
				})
				return
			}
			defer resp.Body.Close()

			limit := int64(maxResponseBodySize)
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, limit))

			el.RunOnLoop(func() {
				res := vm.NewObject()
				_ = res.Set("Status", resp.Status)
				_ = res.Set("StatusCode", resp.StatusCode)
				_ = res.Set("Body", string(respBody))
				resolve(res)
			})
		}()

		return p
	})

	_ = obj.Set("Fetch", func(call sobek.FunctionCall) sobek.Value {
		url := call.Argument(0).String()
		p, resolve, reject := el.CreatePromise()

		go func() {
			resp, err := httpClient.Get(url)
			if err != nil {
				el.RunOnLoop(func() {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch error: %v", err)))
				})
				return
			}
			defer resp.Body.Close()

			limit := int64(maxResponseBodySize)
			body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))

			el.RunOnLoop(func() {
				if err != nil && err != io.EOF {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch body read error: %v", err)))
					return
				}
				if int64(len(body)) > limit {
					reject(vm.NewTypeError(fmt.Sprintf("Fetch response too large (max %d MB)", maxResponseBodySize/1024/1024)))
					return
				}

				res := vm.NewObject()
				_ = res.Set("Status", resp.Status)
				_ = res.Set("StatusCode", resp.StatusCode)
				_ = res.Set("Body", string(body))
				resolve(res)
			})
		}()

		return p
	})

	_ = obj.Set("ListenAndServe", func(call sobek.FunctionCall) sobek.Value {
		addr := call.Argument(0).String()
		handler, ok := sobek.AssertFunction(call.Argument(1))
		if !ok {
			panic(vm.NewTypeError("ListenAndServe requires a handler function"))
		}

		if err := server.ListenAndServe(addr, handler); err != nil {
			panic(vm.NewGoError(err))
		}

		// Return the server for later shutdown
		srvObj := vm.NewObject()
		_ = srvObj.Set("close", func(call sobek.FunctionCall) sobek.Value {
			timeout := 5 * time.Second
			if len(call.Arguments) > 0 {
				timeout = time.Duration(call.Argument(0).ToInteger()) * time.Millisecond
			}
			if err := server.Close(timeout); err != nil {
				panic(vm.NewGoError(err))
			}
			return sobek.Undefined()
		})
		return srvObj
	})

	_ = vm.Set("__go_http__", obj)
}

type stringReader struct {
	s string
	i int
}

func (r *stringReader) Read(b []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += n
	return
}
