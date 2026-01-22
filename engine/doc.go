// Engine wraps the Goja JavaScript runtime with an event loop, memory management,
// and bridge bindings to Go packages. It serves as the main entry point for
// running TypeScript/JavaScript code in the TypeGo runtime.
//
// # Creating an Engine
//
// Use NewEngine to create a fully initialized runtime with all standard bindings:
//
//	eng := engine.NewEngine(128*1024*1024, nil) // 128MB memory limit
//	defer eng.EventLoop.Stop()
//
//	eng.EventLoop.RunOnLoop(func() {
//	    val, err := eng.Run(`console.log("Hello from TypeGo")`)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	})
//
//	eng.EventLoop.Start()
//
// # Memory Management
//
// The memoryLimit parameter sets a soft limit on JavaScript heap size. The engine
// monitors memory usage and can trigger emergency cleanup if limits are exceeded.
//
// # Workers
//
// The engine supports spawning worker threads via the SpawnWorker method. Workers
// run in isolated Goja runtimes but can share memory through the MemoryFactory.
//
// # Event Loop
//
// All JavaScript execution must occur on the event loop. Use RunOnLoop to schedule
// work, and WGAdd/WGDone to track pending async operations. The loop automatically
// stops when all work is complete.
//
// # Error Handling
//
// The engine provides multiple levels of error handling:
//
// 1. JS Errors: Standard JavaScript errors (throw new Error(...)) are returned
// as goja.Exception from Run() and RunSafe().
//
// 2. Go Panics: When using RunSafe(), Go panics are recovered and wrapped as
// Go errors. The original panic value is preserved in the error message.
//
// 3. OnError Callback: Set engine.OnError to receive notifications when errors
// occur in RunSafe(). The callback receives both the error and a stack trace.
//
//	eng := engine.NewEngine(0, nil)
//	eng.OnError = func(err error, stack string) {
//	    log.Printf("Error: %v\nStack:\n%s", err, stack)
//	}
//
// # Context and Cancellation
//
// The engine supports Go context for cancellation:
//
//	ctx := eng.Context()  // Access the engine's context
//	eng.Close()           // Cancels the context and stops the event loop
//
// # Graceful Shutdown
//
// For production use, prefer graceful shutdown over Close():
//
//	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := eng.EventLoop.Shutdown(timeout); err != nil {
//	    log.Printf("Shutdown timed out: %v", err)
//	}
//
// The Shutdown method waits for pending jobs to complete or times out.
package engine
