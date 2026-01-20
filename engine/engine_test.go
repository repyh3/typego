package engine_test

import (
	"testing"
	"time"

	"github.com/repyh3/typego/engine"
)

func TestEngine_Init(t *testing.T) {
	eng := engine.NewEngine(64*1024*1024, nil) // 64MB
	defer eng.Close()

	if eng.VM == nil {
		t.Fatal("Expected VM to be initialized")
	}
	if eng.EventLoop == nil {
		t.Fatal("Expected EventLoop to be initialized")
	}
}

func TestEngine_Run_Basic(t *testing.T) {
	eng := engine.NewEngine(0, nil)
	defer eng.Close()

	val, err := eng.Run(`const a = 1; a + 1`)
	if err != nil {
		t.Fatalf("Engine run failed: %v", err)
	}

	if val.ToInteger() != 2 {
		t.Errorf("Expected 2, got %v", val.ToInteger())
	}
}

func TestEngine_Run_Error(t *testing.T) {
	eng := engine.NewEngine(0, nil)
	defer eng.Close()

	_, err := eng.Run(`throw new Error("test error")`)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// TestEngine_Shutdown verifies clean resource cleanup
func TestEngine_Shutdown(t *testing.T) {
	eng := engine.NewEngine(0, nil)

	// Create a background task
	done := make(chan bool)
	go func() {
		eng.EventLoop.Start()
		done <- true
	}()

	// Allow loop to start
	time.Sleep(10 * time.Millisecond)

	eng.Close()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("EventLoop did not shut down in time")
	}
}
