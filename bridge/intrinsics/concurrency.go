package intrinsics

import (
	"fmt"
	"reflect"

	"github.com/grafana/sobek"
)

// Go implements the typego.go() intrinsic.
func (r *Registry) Go(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		panic(r.vm.NewGoError(newPanicError("go requires a function to execute")))
	}

	fn, ok := sobek.AssertFunction(call.Arguments[0])
	if !ok {
		panic(r.vm.NewGoError(newPanicError("go argument must be a function")))
	}

	args := append([]sobek.Value{}, call.Arguments[1:]...)

	go func() {
		defer func() {
			if rGo := recover(); rGo != nil {
				fmt.Printf("[GOROUTINE PANIC] %v\n", rGo)
			}
		}()

		// Execute the function in the VM
		// NOTE: Sobek is NOT thread-safe for concurrent access to the SAME VM.
		// We use the Registry's VMLock to ensure only one goroutine (or main thread) executes JS at a time.
		r.VMLock.Lock()
		_, err := fn(sobek.Undefined(), args...)
		r.VMLock.Unlock()

		if err != nil {
			fmt.Printf("[GOROUTINE ERROR] %v\n", err)
		}
	}()

	return sobek.Undefined()
}

// Chan represents a bridged Go channel.
type Chan struct {
	ch chan sobek.Value
}

func (c *Chan) Send(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		return sobek.Undefined()
	}
	c.ch <- call.Arguments[0]
	return sobek.Undefined()
}

func (c *Chan) Recv(call sobek.FunctionCall) sobek.Value {
	val, ok := <-c.ch
	if !ok {
		return sobek.Undefined()
	}
	return val
}

func (c *Chan) Close(call sobek.FunctionCall) sobek.Value {
	close(c.ch)
	return sobek.Undefined()
}

// MakeChan implements makeChan(size)
func (r *Registry) MakeChan(call sobek.FunctionCall) sobek.Value {
	size := 0
	if len(call.Arguments) > 0 {
		size = int(call.Arguments[0].ToInteger())
	}

	c := &Chan{ch: make(chan sobek.Value, size)}

	obj := r.vm.NewObject()
	_ = obj.Set("send", c.Send)
	_ = obj.Set("recv", c.Recv)
	_ = obj.Set("close", c.Close)

	// Tag the object so Select can identify the channel
	_ = obj.Set("__chan", c)

	return obj
}

// Select implements select([{ chan, send, recv, default }])
func (r *Registry) Select(call sobek.FunctionCall) sobek.Value {
	if len(call.Arguments) < 1 {
		return sobek.Undefined()
	}

	casesArr := call.Arguments[0].ToObject(r.vm)
	n := int(casesArr.Get("length").ToInteger())

	var selectCases []reflect.SelectCase
	var caseObjs []*sobek.Object

	for i := 0; i < n; i++ {
		caseObj := casesArr.Get(fmt.Sprintf("%d", i)).ToObject(r.vm)
		caseObjs = append(caseObjs, caseObj)

		// 1. Default case
		if defValue := caseObj.Get("default"); defValue != nil && !sobek.IsUndefined(defValue) {
			selectCases = append(selectCases, reflect.SelectCase{
				Dir: reflect.SelectDefault,
			})
			continue
		}

		// 2. Channel operation
		chWrapper := caseObj.Get("chan")
		if chWrapper == nil || sobek.IsUndefined(chWrapper) {
			continue
		}

		chObj := chWrapper.ToObject(r.vm)
		rawChan := chObj.Get("__chan")
		if rawChan == nil {
			continue
		}

		c, ok := rawChan.Export().(*Chan)
		if !ok {
			continue
		}

		chValue := reflect.ValueOf(c.ch)

		// Determine Direction
		if sendVal := caseObj.Get("send"); sendVal != nil && !sobek.IsUndefined(sendVal) {
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: chValue,
				Send: reflect.ValueOf(sendVal),
			})
		} else {
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: chValue,
			})
		}
	}

	if len(selectCases) == 0 {
		return sobek.Undefined()
	}

	// YIELD THE LOCK: Allow goroutines to execute JS while we wait on channels
	r.VMLock.Unlock()
	chosen, recv, recvOk := reflect.Select(selectCases)
	r.VMLock.Lock()

	chosenObj := caseObjs[chosen]

	// Handle Result
	if selectCases[chosen].Dir == reflect.SelectRecv {
		recvVal := recv.Interface().(sobek.Value)
		if recvCb := chosenObj.Get("recv"); recvCb != nil && !sobek.IsUndefined(recvCb) {
			if cb, ok := sobek.AssertFunction(recvCb); ok {
				_, _ = cb(sobek.Undefined(), recvVal, r.vm.ToValue(recvOk))
			}
		}
	} else if selectCases[chosen].Dir == reflect.SelectSend {
		if caseCb := chosenObj.Get("case"); caseCb != nil && !sobek.IsUndefined(caseCb) {
			if cb, ok := sobek.AssertFunction(caseCb); ok {
				_, _ = cb(sobek.Undefined())
			}
		}
	} else if selectCases[chosen].Dir == reflect.SelectDefault {
		if defCb := chosenObj.Get("default"); defCb != nil && !sobek.IsUndefined(defCb) {
			if cb, ok := sobek.AssertFunction(defCb); ok {
				_, _ = cb(sobek.Undefined())
			}
		}
	}

	return r.vm.ToValue(chosen)
}
