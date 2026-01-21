// Package crypto provides bindings for Go's crypto packages.
package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/dop251/goja"
	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/eventloop"
)

func init() {
	core.RegisterModule(&cryptoModule{})
}

type cryptoModule struct{}

func (m *cryptoModule) Name() string {
	return "go:crypto"
}

func (m *cryptoModule) Register(vm *goja.Runtime, el *eventloop.EventLoop) {
	Register(vm)
}

func Register(vm *goja.Runtime) {
	obj := vm.NewObject()

	_ = obj.Set("Sha256", func(call goja.FunctionCall) goja.Value {
		data := call.Argument(0).String()
		hash := sha256.Sum256([]byte(data))
		return vm.ToValue(hex.EncodeToString(hash[:]))
	})

	_ = obj.Set("Sha512", func(call goja.FunctionCall) goja.Value {
		data := call.Argument(0).String()
		hash := sha512.Sum512([]byte(data))
		return vm.ToValue(hex.EncodeToString(hash[:]))
	})

	_ = obj.Set("HmacSha256", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		data := call.Argument(1).String()
		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(data))
		return vm.ToValue(hex.EncodeToString(mac.Sum(nil)))
	})

	_ = obj.Set("HmacSha256Verify", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0).String()
		data := call.Argument(1).String()
		signature := call.Argument(2).String()

		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(data))
		expected := hex.EncodeToString(mac.Sum(nil))

		return vm.ToValue(hmac.Equal([]byte(expected), []byte(signature)))
	})

	_ = obj.Set("RandomBytes", func(call goja.FunctionCall) goja.Value {
		n := int(call.Argument(0).ToInteger())
		if n <= 0 || n > 1024*1024 {
			panic(vm.NewTypeError("RandomBytes: size must be between 1 and 1048576"))
		}
		bytes := make([]byte, n)
		if _, err := rand.Read(bytes); err != nil {
			panic(vm.NewGoError(err))
		}
		return vm.ToValue(hex.EncodeToString(bytes))
	})

	_ = obj.Set("Uuid", func(call goja.FunctionCall) goja.Value {
		uuid := make([]byte, 16)
		if _, err := rand.Read(uuid); err != nil {
			panic(vm.NewGoError(err))
		}

		uuid[6] = (uuid[6] & 0x0f) | 0x40

		uuid[8] = (uuid[8] & 0x3f) | 0x80

		return vm.ToValue(fmt.Sprintf("%x-%x-%x-%x-%x",
			uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]))
	})

	_ = vm.Set("__go_crypto__", obj)
}
