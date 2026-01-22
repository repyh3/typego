// The compiler uses esbuild to bundle TypeScript code, resolve imports, and produce
// a single JavaScript output suitable for execution in the Goja runtime. It handles
// both standard NPM imports and TypeGo's special go: import syntax.
//
// # Basic Usage
//
//	result, err := compiler.Compile("/path/to/main.ts", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.JS) // Bundled JavaScript
//
// # Virtual Modules
//
// The second parameter to Compile allows injecting virtual module content. This is
// used by the TypeGo runtime to provide bindings for go: imports:
//
//	virtualModules := map[string]string{
//	    "go:fmt": `export const Println = (globalThis as any)._go_hyper_fmt.Println;`,
//	}
//	result, err := compiler.Compile(entry, virtualModules)
//
// # Import Resolution
//
// The compiler supports three import namespaces:
//
//   - Standard imports: Resolved from node_modules via esbuild
//   - go/* imports: Internal TypeGo modules (go/fmt, go/os, etc.)
//   - go:* imports: Dynamic Go package bindings (go:fmt, go:github.com/pkg)
//
// The Result struct includes a slice of all go: imports found, enabling the
// runtime to generate bindings before the second compilation pass.
//
// # Two-Pass Compilation
//
// For dynamic go: imports, TypeGo uses a two-pass strategy:
//
//  1. First pass with nil virtualModules to collect import paths
//  2. Generate bindings based on collected imports
//  3. Second pass with populated virtualModules for final bundling
package compiler
