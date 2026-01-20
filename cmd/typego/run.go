package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/repyh3/typego/compiler"
	"github.com/repyh3/typego/internal/builder"
	"github.com/repyh3/typego/internal/linker"
	"github.com/spf13/cobra"
)

var compileMode bool

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a TypeScript file",
	Long: `Run a TypeScript file using the TypeGo engine.

By default, uses interpreter mode for fast execution (<1s startup).
Use --compile to generate a standalone Go binary (slower, but more portable).`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		if compileMode {
			// Standalone compilation mode
			runStandalone(filename)
		} else {
			// Fast interpreter mode (default)
			if err := runInterpreter(filename); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&compileMode, "compile", "c", false, "Compile to standalone binary (slower)")
}

// runStandalone compiles the TypeScript to a standalone Go binary and runs it.
// This is the original behavior, now only used with --compile flag.
func runStandalone(filename string) {
	absPath, _ := filepath.Abs(filename)

	// Create unique temp dir for this run
	tmpDir, err := os.MkdirTemp("", "typego_run_*")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// PASS 1: Scan
	res, _ := compiler.Compile(absPath, nil)

	// Initialize Fetcher
	fetcher, err := linker.NewFetcher()
	if err != nil {
		fmt.Printf("Failed to init fetcher: %v\n", err)
		os.Exit(1)
	}
	defer fetcher.Cleanup()

	// Generate Virtual Modules
	virtualModules := make(map[string]string)
	var bindBlock string

	if res != nil {
		for _, imp := range res.Imports {
			if len(imp) > 3 && imp[:3] == "go:" {
				cleanImp := imp[3:]
				// Skip internal modules handled by bridge
				switch cleanImp {
				case "fmt", "os", "sync", "net/http", "memory":
					continue
				}
				if err := fetcher.Get(cleanImp); err == nil {
					if info, err := linker.Inspect(cleanImp, fetcher.TempDir); err == nil {
						bindBlock += linker.GenerateShim(info, "pkg_"+info.Name)
						var vmContent strings.Builder
						for _, fn := range info.Exports {
							vmContent.WriteString(fmt.Sprintf("export const %s = (globalThis as any)._go_hyper_%s.%s;\n", fn.Name, info.Name, fn.Name))
						}
						virtualModules[imp] = vmContent.String()
					}
				}
			}
		}
	}

	// PASS 2: Compile
	res, err = compiler.Compile(absPath, virtualModules)
	if err != nil {
		fmt.Printf("Build Error: %v\n", err)
		os.Exit(1)
	}

	// Generate Shim
	var importBlock strings.Builder
	for _, imp := range res.Imports {
		if len(imp) > 3 && imp[:3] == "go:" {
			cleanImp := imp[3:]
			// Skip internal modules handled by bridge
			switch cleanImp {
			case "fmt", "os", "sync", "net/http", "memory":
				continue
			}
			importBlock.WriteString(fmt.Sprintf("\t\"%s\"\n", cleanImp))
		}
	}

	// Use shared ShimTemplate from internal/builder
	shimContent := fmt.Sprintf(builder.ShimTemplate, importBlock.String(), fmt.Sprintf("%q", res.JS), bindBlock, memoryLimit*1024*1024)
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(shimContent), 0644); err != nil {
		fmt.Printf("Error writing shim: %v\n", err)
		os.Exit(1)
	}

	// Generate go.mod (minimal, let go get figure out versions)
	goModContent := `module typego_run

go 1.23.6
`
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)

	// Point to local TypeGo source ONLY if we are in the repo (Dev Mode)
	// Otherwise, use the published version
	cwd, _ := os.Getwd()
	absCwd, _ := filepath.Abs(cwd)
	isLocalDev := false
	if data, err := os.ReadFile(filepath.Join(absCwd, "go.mod")); err == nil {
		if strings.Contains(string(data), "module github.com/repyh3/typego") {
			isLocalDev = true
		}
	}

	if isLocalDev {
		fmt.Println("🔧 typego dev mode: using local source replacement")
		replaceCmd := exec.Command("go", "mod", "edit", "-replace", "github.com/repyh3/typego="+absCwd)
		replaceCmd.Dir = tmpDir
		replaceCmd.Run()
	}

	// Fetch dependencies (only root package needed, engine pulls in all subpackages)
	fmt.Print("⏳ Resolving dependencies...")
	getCmd := exec.Command("go", "get", "github.com/repyh3/typego/engine")
	getCmd.Dir = tmpDir
	if isLocalDev {
		getCmd.Env = append(os.Environ(), "GOPROXY=off")
	}
	if err := getCmd.Run(); err != nil {
		fmt.Printf(" failed: %v\n", err)
	} else {
		fmt.Println(" done")
	}

	// Tidy to resolve transitive dependencies
	fmt.Print("⏳ Tidying modules...")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = tmpDir
	if isLocalDev {
		tidy.Env = append(os.Environ(), "GOPROXY=off")
	}
	if err := tidy.Run(); err != nil {
		fmt.Printf(" failed: %v\n", err)
	} else {
		fmt.Println(" done")
	}

	// Build
	fmt.Print("⏳ Compiling...")
	exePath := filepath.Join(tmpDir, "app.exe")
	build := exec.Command("go", "build", "-o", exePath, ".")
	build.Dir = tmpDir
	buildOut, buildErr := build.CombinedOutput()
	if buildErr != nil {
		fmt.Printf(" failed\n%s\n", string(buildOut))
		os.Exit(1)
	}
	fmt.Println(" done")

	// Run
	run := exec.Command(exePath)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if err := run.Run(); err != nil {
		os.Exit(1)
	}
}
