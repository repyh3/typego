// When TypeScript code imports a Go package using the go: prefix (e.g., go:fmt),
// the linker is responsible for fetching, inspecting, and generating bindings
// for that package at compile time.
//
// # Components
//
// The linker consists of three main components:
//
//   - Fetcher: Downloads Go packages using `go get`
//   - Inspector: Analyzes package exports using go/packages
//   - Generator: Produces Go shim code and TypeScript definitions
//
// # Workflow
//
// The typical usage pattern:
//
//	// 1. Create a fetcher
//	fetcher, err := linker.NewFetcher()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer fetcher.Cleanup()
//
//	// 2. Fetch the package
//	if err := fetcher.Get("github.com/fatih/color"); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 3. Inspect exported functions
//	info, err := linker.Inspect("github.com/fatih/color", fetcher.TempDir)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// 4. Generate bindings
//	shimCode := linker.GenerateShim(info, "pkg_color")
//	tsTypes := linker.GenerateTypes(info)
//
// # Package Info
//
// The Inspect function returns a PackageInfo struct containing:
//
//   - Name: The package's short name (e.g., "color")
//   - Path: The full import path (e.g., "github.com/fatih/color")
//   - Exports: A slice of FunctionInfo describing each exported function
//
// # Shim Generation
//
// GenerateShim produces Go code that binds package functions to the Goja runtime.
// The generated code creates a global object (e.g., _go_hyper_color) with methods
// that forward calls to the actual Go functions.
//
// # Type Generation
//
// GenerateTypes produces TypeScript declaration content for IDE support.
// The output includes JSDoc comments derived from the Go source documentation.
package linker
