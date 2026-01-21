// Package http provides HTTP server bindings for TypeGo.
package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/repyh/typego/eventloop"
)

// Server wraps http.Server for TypeGo
type Server struct {
	server *http.Server
	el     *eventloop.EventLoop
	vm     *goja.Runtime
	mu     sync.Mutex
}

// NewServer creates a new TypeGo HTTP server
func NewServer(vm *goja.Runtime, el *eventloop.EventLoop) *Server {
	return &Server{
		vm: vm,
		el: el,
	}
}

// ListenAndServe starts the HTTP server with a JS handler
func (s *Server) ListenAndServe(addr string, handler goja.Callable) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.server = &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create JS Request and Response wrappers
			req := s.wrapRequest(r)
			res := s.wrapResponse(w, r)

			done := make(chan struct{})

			s.el.RunOnLoop(func() {
				defer close(done)
				// Call the JS handler with (req, res)
				_, err := handler(goja.Undefined(), req, res)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			})

			// Wait for handler to complete
			<-done
		}),
	}

	// Start server in background
	s.el.WGAdd(1)
	go func() {
		defer s.el.WGDone()
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return nil
}

// Close shuts down the server gracefully
func (s *Server) Close(timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// wrapRequest creates a JS object representing the HTTP request
func (s *Server) wrapRequest(r *http.Request) goja.Value {
	req := s.vm.NewObject()

	// Basic properties
	_ = req.Set("method", r.Method)
	_ = req.Set("url", r.URL.String())
	_ = req.Set("path", r.URL.Path)
	_ = req.Set("host", r.Host)
	_ = req.Set("proto", r.Proto)

	// Query parameters
	query := s.vm.NewObject()
	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			_ = query.Set(k, v[0])
		} else {
			_ = query.Set(k, v)
		}
	}
	_ = req.Set("query", query)

	// Headers
	headers := s.vm.NewObject()
	for k, v := range r.Header {
		if len(v) == 1 {
			_ = headers.Set(k, v[0])
		} else {
			_ = headers.Set(k, v)
		}
	}
	_ = req.Set("headers", headers)

	// Body reader (async)
	_ = req.Set("body", func(call goja.FunctionCall) goja.Value {
		p, resolve, reject := s.el.CreatePromise()

		go func() {
			body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024)) // 10MB limit
			s.el.RunOnLoop(func() {
				if err != nil {
					reject(s.vm.NewGoError(err))
					return
				}
				resolve(string(body))
			})
		}()

		return p
	})

	// Convenience: get body sync (for small payloads)
	_ = req.Set("bodySync", func(call goja.FunctionCall) goja.Value {
		body, err := io.ReadAll(io.LimitReader(r.Body, 10*1024*1024))
		if err != nil {
			panic(s.vm.NewGoError(err))
		}
		return s.vm.ToValue(string(body))
	})

	return req
}

// wrapResponse creates a JS object for writing HTTP responses
func (s *Server) wrapResponse(w http.ResponseWriter, r *http.Request) goja.Value {
	res := s.vm.NewObject()
	headersSent := false
	statusCode := 200

	_ = res.Set("setHeader", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		value := call.Argument(1).String()
		w.Header().Set(key, value)
		return goja.Undefined()
	})

	_ = res.Set("status", func(call goja.FunctionCall) goja.Value {
		statusCode = int(call.Argument(0).ToInteger())
		return res // Chainable
	})

	_ = res.Set("write", func(call goja.FunctionCall) goja.Value {
		if !headersSent {
			w.WriteHeader(statusCode)
			headersSent = true
		}
		data := call.Argument(0).String()
		_, _ = w.Write([]byte(data))
		return res // Chainable
	})

	// Send response and end
	_ = res.Set("send", func(call goja.FunctionCall) goja.Value {
		if !headersSent {
			w.WriteHeader(statusCode)
			headersSent = true
		}
		if len(call.Arguments) > 0 {
			data := call.Argument(0).String()
			_, _ = w.Write([]byte(data))
		}
		return goja.Undefined()
	})

	// Send JSON response
	_ = res.Set("json", func(call goja.FunctionCall) goja.Value {
		w.Header().Set("Content-Type", "application/json")
		if !headersSent {
			w.WriteHeader(statusCode)
			headersSent = true
		}
		if len(call.Arguments) > 0 {
			val := call.Argument(0).Export()
			jsonStr, err := toJSON(val)
			if err != nil {
				panic(s.vm.NewGoError(err))
			}
			_, _ = w.Write([]byte(jsonStr))
		}
		return goja.Undefined()
	})

	// Redirect
	_ = res.Set("redirect", func(call goja.FunctionCall) goja.Value {
		url := call.Argument(0).String()
		code := 302
		if len(call.Arguments) > 1 {
			code = int(call.Argument(1).ToInteger())
		}
		http.Redirect(w, r, url, code)
		return goja.Undefined()
	})

	return res
}
