package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/repyh/typego/internal/ecosystem"
	"github.com/repyh/typego/internal/ecosystem/ast"
	"github.com/repyh/typego/internal/ecosystem/builder"
	"github.com/repyh/typego/internal/ecosystem/resolver"
	"github.com/repyh/typego/internal/linker"
)

const (
	ConfigFileName = "typego.modules.json"
	HiddenDirName  = ".typego"
	BinaryName     = "typego-app.exe"
)

// RunInstall executes the installation process
func RunInstall(cwd string) error {
	// 1. Read Config
	configPath := filepath.Join(cwd, ConfigFileName)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found. Run 'typego init' first", ConfigFileName)
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	var config ecosystem.ModuleConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("invalid config JSON: %w", err)
	}
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 2. Setup Hidden Directory
	hiddenDir := filepath.Join(cwd, HiddenDirName)
	if _, err := os.Stat(hiddenDir); os.IsNotExist(err) {
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", HiddenDirName, err)
		}
	}

	// 3. Scan for Usage (AST Analysis)
	// For now, let's scan src/*.ts
	// In reality we should walk the tree.
	files, _ := filepath.Glob(filepath.Join(cwd, "src", "*.ts"))
	var allUsedModules []string

	// Also check index.ts in root if it exists
	if _, err := os.Stat(filepath.Join(cwd, "index.ts")); err == nil {
		files = append(files, filepath.Join(cwd, "index.ts"))
	}

	fmt.Printf("🔍 Scanning %d files for go:* imports...\n", len(files))

	for _, file := range files {
		res, _, err := ast.ParseFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file, err)
			continue
		}
		// Extract imports from Metafile (ParseFile returns BuildResult which has Metafile)
		// Wait, ast.ParseFile returns *api.BuildResult.
		// We need to pass res.Metafile to Extractor.
		imports, _ := ast.ExtractImports(res.Metafile)
		allUsedModules = append(allUsedModules, imports...)
	}

	// Dedup
	usedMap := make(map[string]bool)
	for _, m := range allUsedModules {
		// remove "go:" prefix logic should be in Extractor or here
		// Extractor returns full "go:github.com/..."
		usedMap[m] = true
	}

	// 4. Resolve Dependencies
	// Initialize Fetcher to inspect types
	fetcher, err := linker.NewFetcher()
	if err != nil {
		return fmt.Errorf("failed to init fetcher: %w", err)
	}
	defer fetcher.Cleanup()

	tsShims := make(map[string]string)
	var bridgeBlock strings.Builder
	namedImports := make(map[string]string)

	fmt.Println("📦 Resolving dependencies...")

	// Create a temp workspace for toolchain operations
	// We use .typego/work
	workDir := filepath.Join(hiddenDir, "work")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}

	// Initialize go.mod in workDir
	_ = resolver.InitGoMod(workDir, "typego_jit_build")

	// Detect if we are in the typego repo (local dev) to add source replacement
	typegoRoot, isLocalDev := ecosystem.FindRepoRoot(cwd)

	if isLocalDev {
		fmt.Printf("🔧 typego dev mode: using local source replacement at %s\n", typegoRoot)
		replaceCmd := exec.Command("go", "mod", "edit", "-replace", "github.com/repyh/typego="+typegoRoot)
		replaceCmd.Dir = workDir
		_ = replaceCmd.Run()
	}

	requireCmd := exec.Command("go", "mod", "edit", "-require", "github.com/repyh/typego@v0.0.0")
	requireCmd.Dir = workDir
	_ = requireCmd.Run()

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = workDir
	_ = tidyCmd.Run()

	// Filter imports that are external
	for m := range usedMap {
		if len(m) > 3 && m[:3] == "go:" {
			cleanName := m[3:]

			// Skip internal
			switch cleanName {
			case "fmt", "os", "sync", "net/http", "memory", "crypto":
				continue
			}

			// Check if allowed in config
			// config.Dependencies keys should match
			// We iterate config to find matches or just try to allow any if config is empty?
			// Strict mode: Only allow if in config.
			// Loose mode: Allow all found.

			// For MVP: Check config. If config has it, use version from config.
			if _, ok := config.Dependencies[cleanName]; !ok {
				fmt.Printf("⚠️  Warning: %s is used but not listed in %s\n", cleanName, ConfigFileName)
				// We can still try to install it with "latest"
			}

			// Go Get
			// If we have version in config, use it.
			version := "latest"
			if v, ok := config.Dependencies[cleanName]; ok {
				version = v
			}

			// run go get in workDir
			fmt.Printf("   Getting %s@%s...\n", cleanName, version)
			if err := resolver.RunGoGet(workDir, []string{cleanName + "@" + version}); err != nil {
				return fmt.Errorf("failed to get %s: %w", cleanName, err)
			}

			// Inspect for Shim Generation
			// We can inspect in workDir!
			if info, err := linker.Inspect(cleanName, workDir); err == nil {
				// Generate TS Shim for Compiler
				tsShims[m] = linker.GenerateTSShim(info)

				// Generate Go Bridge for Engine
				bridgeBlock.WriteString(linker.GenerateShim(info, "pkg_"+info.Name))

				namedImports[cleanName] = info.Name
			} else {
				fmt.Printf("   Failed to inspect %s: %v\n", cleanName, err)
			}
		}
	}

	// Also get typego/engine and cmd/typego in workDir
	// resolver.RunGoGet(workDir, "github.com/repyh/typego@latest") // or replace if dev

	// 5. Scaffold main.go
	fmt.Println("🏗️  Scaffolding binary...")
	if err := builder.ScaffoldMain(workDir, namedImports, tsShims, bridgeBlock.String()); err != nil {
		return err
	}

	// 6. Compile
	fmt.Println("🔨 Compiling JIT binary...")
	binDir := filepath.Join(hiddenDir, "bin")
	if err := builder.CompileBinary(workDir, binDir, BinaryName); err != nil {
		return err
	}

	binaryPath := filepath.Join(binDir, BinaryName)
	fmt.Println("✅ Install Complete! Binary cached at", binaryPath)

	// Ensure .typego is ignored
	_ = ecosystem.EnsureGitIgnore(cwd)

	// Write checksum for stale detection
	if hash, err := ecosystem.GetConfigHash(cwd); err == nil {
		_ = ecosystem.WriteChecksum(cwd, hash)
	}

	return nil
}
