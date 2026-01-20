package integration

import (
	"testing"
	"time"
)

func TestBridge_Fmt(t *testing.T) {
	harness := NewHarness(t)
	// Use internal global __go_fmt__
	harness.Run(t, `
		const fmt = __go_fmt__;
		fmt.Println("Testing Println");
		fmt.Printf("Testing %s\n", "Printf");
	`)
}

func TestBridge_Sync_Sleep(t *testing.T) {
	harness := NewHarness(t)
	start := time.Now()
	// Use internal global __go_sync__
	harness.Run(t, `
		const sync = __go_sync__;
		await sync.Sleep(200);
	`)
	elapsed := time.Since(start)
	if elapsed < 200*time.Millisecond {
		t.Errorf("Sleep duration too short: %v", elapsed)
	}
}

func TestBridge_Worker_Init(t *testing.T) {
	harness := NewHarness(t)
	// Use internal global __typego_worker__
	harness.Run(t, `
		const workerLib = __typego_worker__;
        if (typeof workerLib.Worker !== 'function') {
            throw new Error("Worker is not a function");
        }
        
        // We can't easily spawn a worker without a file on disk,
        // but verifying the constructor exists proves the bridge registered it.
	`)
}
