package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/repyh/typego/bridge/core"
	"github.com/repyh/typego/bridge/intrinsics"
	bridge_crypto "github.com/repyh/typego/bridge/modules/crypto"

	bridge_net "github.com/repyh/typego/bridge/modules/net"
	bridge_sync "github.com/repyh/typego/bridge/modules/sync"
	bridge_memory "github.com/repyh/typego/bridge/stdlib/memory"
	bridge_worker "github.com/repyh/typego/bridge/stdlib/worker"
	"github.com/repyh/typego/compiler"
	"github.com/repyh/typego/internal/linker"
	"github.com/spf13/cobra"
)

var TypesCmd = &cobra.Command{
	Use:   "types [file]",
	Short: "Sync and update TypeGo ambient definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing TypeGo definitions...")

		dtsPath := filepath.Join(".typego", "types", "go.d.ts")
		if err := os.MkdirAll(filepath.Dir(dtsPath), 0755); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		var currentContent []byte
		if existing, err := os.ReadFile(dtsPath); err == nil && len(existing) > 0 {
			currentContent = existing
		}

		// Always ensure global types are up to date at the start
		newGlobal := string(core.GlobalTypes) + "\n" + intrinsics.IntrinsicTypes

		reGlobal := regexp.MustCompile(`(?s)^// TypeGo Type Definitions.*?// TypeGo Namespaces`)
		if reGlobal.Match(currentContent) {
			currentContent = reGlobal.ReplaceAll(currentContent, []byte(newGlobal))
		} else if len(currentContent) == 0 {
			currentContent = []byte(newGlobal)
		} else {
			currentContent = append([]byte(newGlobal+"\n"), currentContent...)
		}

		fetcher, err := linker.NewFetcher()
		if err != nil {
			fmt.Printf("Failed to init fetcher: %v\n", err)
			return
		}
		defer fetcher.Cleanup()

		fileSet := make(map[string]bool)
		if len(args) > 0 {
			absPath, _ := filepath.Abs(args[0])
			fileSet[absPath] = true
		} else {
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

		// Generate types for each unique import (Recursive Resolution)
		processed := make(map[string]bool)
		queue := []string{}
		for imp := range goImports {
			queue = append(queue, imp)
		}

		for len(queue) > 0 {
			imp := queue[0]
			queue = queue[1:]

			if processed[imp] {
				continue
			}
			processed[imp] = true

			// Skip modules that are already fully defined in std.d.ts
			if imp == "go:net/http" || imp == "go:sync" || imp == "typego:memory" || imp == "typego:worker" || imp == "go:memory" || imp == "go:crypto" {
				continue
			}

			fmt.Printf("📦 Generating types for %s...\n", imp)

			var pkgPath string
			if strings.HasPrefix(imp, "go:") {
				pkgPath = strings.TrimPrefix(imp, "go:")
			} else {
				continue
			}

			// specific internal modules don't need fetching if we are in the repo
			isInternal := strings.HasPrefix(pkgPath, "github.com/repyh/typego")
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

			// Discovery Discovery: Add imports from fields/embeds to the queue
			for _, st := range info.Structs {
				for _, field := range st.Fields {
					if field.ImportPath != "" && field.ImportPath != pkgPath {
						newImp := "go:" + field.ImportPath
						if !processed[newImp] {
							queue = append(queue, newImp)
						}
					}
				}
				for _, embed := range st.Embeds {
					if embed.ImportPath != "" && embed.ImportPath != pkgPath {
						newImp := "go:" + embed.ImportPath
						if !processed[newImp] {
							queue = append(queue, newImp)
						}
					}
				}
			}

			// Update ImportPath to match the requested import (e.g. "go:fmt" -> "fmt" so generator adds "go:")
			info.ImportPath = pkgPath

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
