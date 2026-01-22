package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	"github.com/repyh/typego/compiler"
	"github.com/repyh/typego/internal/builder"
	"github.com/repyh/typego/internal/ecosystem"
	"github.com/repyh/typego/internal/linker"
	"github.com/spf13/cobra"
)

var buildOut string
var minify bool
var buildTarget string

var supportedTargets = map[string]struct{ goos, goarch string }{
	"linux-amd64":   {"linux", "amd64"},
	"linux-arm64":   {"linux", "arm64"},
	"darwin-amd64":  {"darwin", "amd64"},
	"darwin-arm64":  {"darwin", "arm64"},
	"windows-amd64": {"windows", "amd64"},
}

var BuildCmd = &cobra.Command{
	Use:   "build [file]",
	Short: "Build and bundle a TypeScript file for production",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		absPath, _ := filepath.Abs(filename)

		fmt.Printf("📦 Building %s...\n", absPath)

		// We ignore errors here because the virtual modules are not yet populated
		res, _ := compiler.Compile(absPath, nil)

		tmpDir := ".typego_build_tmp"
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)

		fetcher, err := linker.NewFetcher()
		if err != nil {
			fmt.Printf("Failed to init fetcher: %v\n", err)
			os.Exit(1)
		}
		defer fetcher.Cleanup()

		virtualModules := make(map[string]string)
		var bindBlock string

		if res != nil {
			for _, imp := range res.Imports {
				if len(imp) > 3 && imp[:3] == "go:" {
					cleanImp := imp[3:]

					fmt.Printf("🔍 Inspecting %s...\n", cleanImp)
					_ = fetcher.Get(cleanImp)

					info, err := linker.Inspect(cleanImp, fetcher.TempDir)
					if err != nil {
						fmt.Printf("Warning: Could not inspect %s: %v\n", cleanImp, err)
						continue
					}

					bindBlock += linker.GenerateShim(info, "pkg_"+info.Name)

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

		fmt.Println("🔨 Compiling binary (Pass 2)...")
		res, err = compiler.Compile(absPath, virtualModules)
		if err != nil {
			fmt.Printf("Build Error: %v\n", err)
			os.Exit(1)
		}

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

		shimContent := fmt.Sprintf(builder.ShimTemplate, importBlock.String(), fmt.Sprintf("%q", res.JS), bindBlock, MemoryLimit*1024*1024)

		shimPath := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(shimPath, []byte(shimContent), 0644); err != nil {
			fmt.Printf("Error writing shim: %v\n", err)
			os.Exit(1)
		}

		goModContent := `module typego_app

go 1.23.6
`

		if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
			fmt.Printf("Error writing go.mod: %v\n", err)
			os.Exit(1)
		}

		cwd, _ := os.Getwd()
		typegoRoot, isLocalDev := ecosystem.FindRepoRoot(cwd)

		if isLocalDev {
			fmt.Println("🔧 typego dev mode: using local source replacement at", typegoRoot)
			replaceCmd := exec.Command("go", "mod", "edit", "-replace", "github.com/repyh/typego="+typegoRoot)
			replaceCmd.Dir = tmpDir
			_ = replaceCmd.Run()
		}

		packages := []string{
			"github.com/repyh/typego/bridge/core",
			"github.com/repyh/typego/bridge/polyfills",
			"github.com/repyh/typego/engine",
			"github.com/repyh/typego/eventloop",
			"github.com/grafana/sobek",
		}
		for _, pkg := range packages {
			// If local dev, we don't want @latest for our own packages, but they are replaced anyway.
			// However, go get might complain if version is missing.
			// Let's just use go get without version, let go.mod handle it via replacement.
			target := pkg
			if !isLocalDev {
				target = pkg + "@latest"
			}
			getCmd := exec.Command("go", "get", target)
			getCmd.Dir = tmpDir
			if isLocalDev {
				getCmd.Env = append(os.Environ(), "GOPROXY=off") // Force local
			} else {
				getCmd.Env = append(os.Environ(), "GOPROXY=direct")
			}
			_ = getCmd.Run()
		}

		outputName := buildOut
		if outputName == "" {
			outputName = "app.exe"
		}
		// Make output absolute so go build puts it in the right place
		absOut, _ := filepath.Abs(outputName)

		fmt.Println("🧹 Resolving dependencies...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = tmpDir
		tidyCmd.Env = append(os.Environ(), "GOPROXY=direct")
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

		env := os.Environ()
		if buildTarget != "" {
			if target, ok := supportedTargets[buildTarget]; ok {
				fmt.Printf("🎯 Targeting %s/%s...\n", target.goos, target.goarch)
				env = append(env, "GOOS="+target.goos, "GOARCH="+target.goarch)
				// Ensure CGO is disabled for cross-compilation unless explicitly needed
				env = append(env, "CGO_ENABLED=0")
			} else {
				fmt.Printf("Error: Unsupported target '%s'. Available targets:\n", buildTarget)
				for t := range supportedTargets {
					fmt.Printf("  - %s\n", t)
				}
				os.Exit(1)
			}
		}
		buildCmd.Env = env

		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Compilation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✨ Binary created: %s\n", outputName)
	},
}

func init() {
	BuildCmd.Flags().StringVarP(&buildOut, "out", "o", "dist/index.js", "Output bundle path")
	BuildCmd.Flags().BoolVarP(&minify, "minify", "m", false, "Minify output")
	BuildCmd.Flags().StringVarP(&buildTarget, "target", "t", "", "Cross-compilation target (e.g. linux-amd64)")
	// Registered in root.go
}
