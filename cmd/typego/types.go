package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/repyh3/typego/bridge/core"
	bridge_crypto "github.com/repyh3/typego/bridge/modules/crypto"
	bridge_net "github.com/repyh3/typego/bridge/modules/net"
	bridge_sync "github.com/repyh3/typego/bridge/modules/sync"
	bridge_memory "github.com/repyh3/typego/bridge/stdlib/memory"
	bridge_worker "github.com/repyh3/typego/bridge/stdlib/worker"
	"github.com/repyh3/typego/compiler"
	"github.com/repyh3/typego/internal/linker"
	"github.com/spf13/cobra"
)

var typesCmd = &cobra.Command{
	Use:   "types [file]",
	Short: "Sync and update TypeGo ambient definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing TypeGo definitions...")

		dtsPath := filepath.Join(".typego", "types", "go.d.ts")
		if err := os.MkdirAll(filepath.Dir(dtsPath), 0755); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Read existing (target) file first, or use embedded default
		var currentContent []byte
		if existing, err := os.ReadFile(dtsPath); err == nil && len(existing) > 0 {
			currentContent = existing
		} else {
			currentContent = core.GlobalTypes
		}

		// Initialize Fetcher
		fetcher, err := linker.NewFetcher()
		if err != nil {
			fmt.Printf("Failed to init fetcher: %v\n", err)
			return
		}
		defer fetcher.Cleanup()

		// Collect all TypeScript files to scan (deduplicated)
		fileSet := make(map[string]bool)
		if len(args) > 0 {
			// Explicit file provided
			absPath, _ := filepath.Abs(args[0])
			fileSet[absPath] = true
		} else {
			// Recursively find all .ts files in the current project
			err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				// Skip common non-source and system directories
				if d.IsDir() {
					name := d.Name()
					if name == "node_modules" || name == ".typego" || name == ".git" || name == ".gemini" {
						return filepath.SkipDir
					}
					return nil
				}

				// Only collect .ts files
				if filepath.Ext(path) == ".ts" {
					absPath, _ := filepath.Abs(path)
					fileSet[absPath] = true
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error searching for TypeScript files: %v\n", err)
			}
		}

		// Collect unique Go imports from all files
		goImports := make(map[string]bool)

		// FORCE: Always include core modules that need inspection
		coreModules := []string{
			"go:fmt", "go:os", "go:net/url",
		}
		for _, mod := range coreModules {
			goImports[mod] = true
		}

		for file := range fileSet {
			fmt.Printf("🔍 Scanning %s...\n", filepath.Base(file))
			res, _ := compiler.Compile(file, nil)
			if res != nil {
				for _, imp := range res.Imports {
					if strings.HasPrefix(imp, "go:") || strings.HasPrefix(imp, "typego:") {
						goImports[imp] = true
					}
				}
			}
		}

		// Generate types for each unique import
		for imp := range goImports {
			// Skip modules that are already fully defined in std.d.ts
			if imp == "go:net/http" || imp == "go:sync" || imp == "typego:memory" || imp == "typego:worker" || imp == "go:memory" || imp == "go:crypto" {
				continue
			}

			fmt.Printf("📦 Generating types for %s...\n", imp)

			var pkgPath string

			if strings.HasPrefix(imp, "go:") {
				name := imp[3:]
				switch name {
				case "fmt", "os", "net/url":
					// Inspect native Go package for high-quality types
					pkgPath = name
				default:
					pkgPath = name
				}
			} else if strings.HasPrefix(imp, "typego:") {
				// typego: modules usually handled in std.d.ts, but if any dynamic ones added later...
				continue
			}

			if pkgPath == "" {
				continue
			}

			// specific internal modules don't need fetching if we are in the repo
			isInternal := strings.HasPrefix(pkgPath, "github.com/repyh3/typego")
			if !isInternal {
				if err := fetcher.Get(pkgPath); err != nil {
					fmt.Printf("Warning: Failed to fetch %s: %v\n", pkgPath, err)
					continue
				}
			}

			info, err := linker.Inspect(pkgPath, fetcher.TempDir)
			if err != nil {
				fmt.Printf("Failed to inspect %s: %v\n", pkgPath, err)
				continue
			}

			// Update ImportPath to match the requested import (e.g. "go:fmt" -> "fmt" so generator adds "go:")
			info.ImportPath = strings.TrimPrefix(imp, "go:")

			newTypeBlock := linker.GenerateTypes(info)

			// Additive Merge Logic
			pattern := fmt.Sprintf(`(?s)// MODULE: %s.*?// END: %s\n`, regexp.QuoteMeta(imp), regexp.QuoteMeta(imp))
			re := regexp.MustCompile(pattern)

			if re.Match(currentContent) {
				currentContent = re.ReplaceAll(currentContent, []byte(newTypeBlock))
			} else {
				currentContent = append(currentContent, []byte("\n"+newTypeBlock)...)
			}
		}

		currentContent = updateTypeBlock(currentContent, "go:net/http", string(bridge_net.HttpTypes))
		currentContent = updateTypeBlock(currentContent, "go:sync", string(bridge_sync.Types))
		currentContent = updateTypeBlock(currentContent, "go:memory", string(bridge_memory.Types))
		currentContent = updateTypeBlock(currentContent, "typego:memory", string(bridge_memory.Types))
		currentContent = updateTypeBlock(currentContent, "typego:worker", string(bridge_worker.Types))
		currentContent = updateTypeBlock(currentContent, "go:crypto", string(bridge_crypto.Types))

		if err := os.WriteFile(dtsPath, currentContent, 0644); err != nil {
			fmt.Printf("Error writing types: %v\n", err)
			return
		}

		fmt.Println("✅ Definitions synced to .typego/types/go.d.ts")
	},
}

func init() {
	rootCmd.AddCommand(typesCmd)
}

func updateTypeBlock(content []byte, moduleName, typeDef string) []byte {
	block := typeDef

	pattern := fmt.Sprintf(`(?s)// MODULE: %s.*?// END: %s\n`, regexp.QuoteMeta(moduleName), regexp.QuoteMeta(moduleName))
	re := regexp.MustCompile(pattern)

	if re.Match(content) {
		return re.ReplaceAll(content, []byte(block))
	}
	return append(content, []byte("\n"+block)...)
}
