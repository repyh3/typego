package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	"github.com/repyh3/typego/compiler"
	"github.com/repyh3/typego/internal/linker"
	"github.com/spf13/cobra"
)

var buildOut string
var minify bool

var buildCmd = &cobra.Command{
	Use:   "build [file]",
	Short: "Build and bundle a TypeScript file for production",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		absPath, _ := filepath.Abs(filename)

		fmt.Printf("📦 Building %s...\n", absPath)

		// PASS 1: Scan for imports
		// We ignore errors here because the virtual modules are not yet populated
		res, _ := compiler.Compile(absPath, nil)

		// 2. Prepare Temp Directory
		tmpDir := ".typego_build_tmp"
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir) // Cleanup

		// Initialize Fetcher
		fetcher, err := linker.NewFetcher()
		if err != nil {
			fmt.Printf("Failed to init fetcher: %v\n", err)
			os.Exit(1)
		}
		defer fetcher.Cleanup()

		// 3. Inspect Imports & Generate Virtual Modules
		virtualModules := make(map[string]string)
		var bindBlock string

		if res != nil {
			for _, imp := range res.Imports {
				if len(imp) > 3 && imp[:3] == "go:" {
					cleanImp := imp[3:]

					fmt.Printf("🔍 Inspecting %s...\n", cleanImp)
					if err := fetcher.Get(cleanImp); err != nil {
						// Fallback: It might be stdlib, ignore error for now
					}

					info, err := linker.Inspect(cleanImp, fetcher.TempDir)
					if err != nil {
						fmt.Printf("Warning: Could not inspect %s: %v\n", cleanImp, err)
						continue
					}

					// Generate Shim Binding Code
					bindBlock += linker.GenerateShim(info, "pkg_"+info.Name)

					// Generate Virtual Module Content (Exports)
					var vmContent strings.Builder
					for _, fn := range info.Exports {
						// e.g. export const Red = (globalThis as any)._go_hyper_color.Red;
						// Note: variable name matches GenerateShim output: "_go_hyper_" + info.Name
						vmContent.WriteString(fmt.Sprintf("export const %s = (globalThis as any)._go_hyper_%s.%s;\n", fn.Name, info.Name, fn.Name))
					}

					virtualModules[imp] = vmContent.String()
				}
			}
		}

		// PASS 2: Compile with Virtual Modules (Strict)
		fmt.Println("🔨 Compiling binary (Pass 2)...")
		res, err = compiler.Compile(absPath, virtualModules)
		if err != nil {
			fmt.Printf("Build Error: %v\n", err)
			os.Exit(1)
		}

		// Generate Import Block for Shim
		var importBlock strings.Builder
		for _, imp := range res.Imports {
			if len(imp) > 3 && imp[:3] == "go:" {
				cleanImp := imp[3:]
				// Skip hardcoded imports in shimTemplate
				if cleanImp == "fmt" || cleanImp == "os" {
					continue
				}
				importBlock.WriteString(fmt.Sprintf("\t\"%s\"\n", cleanImp))
			}
		}

		shimContent := fmt.Sprintf(shimTemplate, importBlock.String(), fmt.Sprintf("%q", res.JS), bindBlock)

		shimPath := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err != nil {
			fmt.Printf("Error writing shim: %v\n", err)
			os.Exit(1)
		}

		// 4. Generate go.mod for the shim
		currWd, _ := os.Getwd()
		goModContent := fmt.Sprintf(`module typego_app

go 1.23.6

require github.com/repyh3/typego v0.0.0

replace github.com/repyh3/typego => %s
`, filepath.ToSlash(currWd))

		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
			fmt.Printf("Error writing go.mod: %v\n", err)
			os.Exit(1)
		}

		// 5. Build Binary
		outputName := buildOut
		if outputName == "" {
			outputName = "app.exe"
		}
		// Make output absolute so go build puts it in the right place
		absOut, _ := filepath.Abs(outputName)

		fmt.Println("🧹 Resolving dependencies...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = tmpDir
		tidyCmd.Stdout = os.Stdout
		tidyCmd.Stderr = os.Stderr
		if err := tidyCmd.Run(); err != nil {
			fmt.Printf("go mod tidy failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🔨 Compiling binary...")
		buildCmd := exec.Command("go", "build", "-o", absOut, ".")
		buildCmd.Dir = tmpDir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Compilation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✨ Binary created: %s\n", outputName)
	},
}

const shimTemplate = `package main

import (
	"fmt"
	"os"
	"time"

	%s

	"github.com/dop251/goja"
	"github.com/repyh3/typego/bridge"
	"github.com/repyh3/typego/engine"
)

const jsBundle = %s

type NativeTools struct {
	StartTime string
}

func (n *NativeTools) GetRuntimeInfo() string {
	return "TypeGo Standalone v1.0"
}

func main() {
	eng := engine.NewEngine(128*1024*1024, nil)

	// Initialize Shared Buffer
	cliBuffer := make([]byte, 1024)
	bridge.MapSharedBuffer(eng.VM, "cliBuffer", cliBuffer)

	// Initialize Memory Factory (for makeShared)
	sharedBuffers := make(map[string][]byte)
	bridge.EnableMemoryFactory(eng.VM, sharedBuffers)

	// Initialize Worker API
	bridge.EnableWorkerAPI(eng.VM, eng.EventLoop)

	// Initialize Native Tools
	tools := &NativeTools{StartTime: "2026-01-16"}
	_ = bridge.BindStruct(eng.VM, "native", tools)

	// Node.js Polyfills (Process & Buffer)
	bridge.EnableNodeCompat(eng.VM)
	_, _ = eng.VM.RunString("if (typeof Buffer === 'undefined') { globalThis.Buffer = { from: function(str) { if (typeof str === 'string') { var arr = []; for (var i = 0; i < str.length; i++) { arr.push(str.charCodeAt(i)); } return new Uint8Array(arr); } return new Uint8Array(str); }, alloc: function(size) { return new Uint8Array(size); } }; }")

	// Timer Polyfill
	eng.VM.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		fn, _ := goja.AssertFunction(call.Argument(0))
		ms := call.Argument(1).ToInteger()
		eng.EventLoop.WGAdd(1)
		go func() {
			time.Sleep(time.Duration(ms) * time.Millisecond)
			eng.EventLoop.RunOnLoop(func() {
				_, _ = fn(goja.Undefined())
				eng.EventLoop.WGDone()
			})
		}()
		return goja.Undefined()
	})

	// Hyper-Linker Bindings (Generated)
	%s

	// Run on EventLoop
	eng.EventLoop.RunOnLoop(func() {
		val, err := eng.Run(jsBundle)
		if err != nil {
			fmt.Printf("Runtime Error: %%v\n", err)
			os.Exit(1)
		}

		// Handle Top-Level Async (Promises)
		if val != nil && !goja.IsUndefined(val) && !goja.IsNull(val) {
			if obj := val.ToObject(eng.VM); obj != nil {
				then := obj.Get("then")
				if then != nil && !goja.IsUndefined(then) {
					if _, ok := goja.AssertFunction(then); ok {
						eng.EventLoop.WGAdd(1)
						done := eng.VM.ToValue(func(goja.FunctionCall) goja.Value {
							eng.EventLoop.WGDone()
							return goja.Undefined()
						})
						thenFn, _ := goja.AssertFunction(then)
						_, _ = thenFn(val, done, done)
					}
				}
			}
		}
	})

	eng.EventLoop.Start()
}
`

func init() {
	buildCmd.Flags().StringVarP(&buildOut, "out", "o", "dist/index.js", "Output bundle path")
	buildCmd.Flags().BoolVarP(&minify, "minify", "m", false, "Minify output")
	rootCmd.AddCommand(buildCmd)
}
