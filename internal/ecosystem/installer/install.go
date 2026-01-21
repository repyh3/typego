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

func RunInstall(cwd string) error {
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

	hiddenDir := filepath.Join(cwd, HiddenDirName)
	if _, err := os.Stat(hiddenDir); os.IsNotExist(err) {
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", HiddenDirName, err)
		}
	}

	// Scan for go:* imports in TypeScript files
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

	fetcher, err := linker.NewFetcher()
	if err != nil {
		return fmt.Errorf("failed to init fetcher: %w", err)
	}
	defer fetcher.Cleanup()

	tsShims := make(map[string]string)
	var bridgeBlock strings.Builder
	namedImports := make(map[string]string)
	var packageInfos []*linker.PackageInfo // For auto-types generation

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

	requireCmd := exec.Command("go", "mod", "edit", "-require", "github.com/repyh/typego@v1.3.1")
	requireCmd.Dir = workDir
	_ = requireCmd.Run()

	// Process external go:* imports
	for m := range usedMap {
		if len(m) > 3 && m[:3] == "go:" {
			cleanName := m[3:]

			// Skip internal/stdlib packages
			switch cleanName {
			case "fmt", "os", "sync", "net/http", "memory", "crypto":
				continue
			}

			// Check if in config, warn if not
			if _, ok := config.Dependencies[cleanName]; !ok {
				fmt.Printf("⚠️  Warning: %s is used but not in %s\n", cleanName, ConfigFileName)
			}

			// Determine version
			version := "latest"
			if v, ok := config.Dependencies[cleanName]; ok {
				version = v
			}

			// Fetch the package
			fmt.Printf("   📥 Getting %s@%s...\n", cleanName, version)
			if err := resolver.RunGoGet(workDir, []string{cleanName + "@" + version}); err != nil {
				return fmt.Errorf("failed to get %s: %w", cleanName, err)
			}

			// Inspect for shim generation
			if info, err := linker.Inspect(cleanName, workDir); err == nil {
				tsShims[m] = linker.GenerateTSShim(info)
				bridgeBlock.WriteString(linker.GenerateShim(info, "pkg_"+info.Name))
				namedImports[cleanName] = info.Name
				packageInfos = append(packageInfos, info) // Collect for types
			} else {
				fmt.Printf("   ⚠️  Failed to inspect %s: %v\n", cleanName, err)
			}
		}
	}

	// Scaffold main.go before tidy so go mod tidy sees the imports
	fmt.Println("🏗️  Scaffolding binary...")
	if err := builder.ScaffoldMain(workDir, namedImports, tsShims, bridgeBlock.String()); err != nil {
		return err
	}

	fmt.Println("🧹 Resolving dependencies...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = workDir
	// Allow transitive dependency resolution even with local replacement
	if err := tidyCmd.Run(); err != nil {
		// Output the error from tidy if it fails
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("⚠️ go mod tidy failed: %s\n", string(exitErr.Stderr))
		} else {
			fmt.Printf("⚠️ go mod tidy failed: %v\n", err)
		}
	}

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

	// Generate lockfile from resolved go.mod
	writeLockfile(cwd, workDir, config.Dependencies)

	// Auto-generate types
	if len(packageInfos) > 0 {
		fmt.Println("📦 Syncing type definitions...")
		writeTypeDefinitions(cwd, packageInfos)
	}

	return nil
}

func writeLockfile(projectRoot, workDir string, deps map[string]string) {
	goModPath := filepath.Join(workDir, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return
	}

	lockfile := ecosystem.DefaultLockfile()

	// Parse require directives to extract resolved versions
	lines := strings.Split(string(data), "\n")
	inRequire := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if inRequire && len(line) > 0 && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				pkgPath := parts[0]
				version := parts[1]
				// Only include user dependencies, not transitive
				if _, isDep := deps[pkgPath]; isDep {
					lockfile.Resolved[pkgPath] = ecosystem.ResolvedPackage{Version: version}
				}
			}
		}
	}

	if len(lockfile.Resolved) == 0 {
		return
	}

	lockData, _ := json.MarshalIndent(lockfile, "", "  ")
	lockPath := filepath.Join(projectRoot, ecosystem.LockFileName)
	_ = os.WriteFile(lockPath, lockData, 0644)
	fmt.Printf("🔒 Lockfile written to %s\n", ecosystem.LockFileName)
}

func writeTypeDefinitions(projectRoot string, infos []*linker.PackageInfo) {
	typesDir := filepath.Join(projectRoot, ".typego", "types")
	_ = os.MkdirAll(typesDir, 0755)

	var sb strings.Builder
	for _, info := range infos {
		sb.WriteString(linker.GenerateTypes(info))
		sb.WriteString("\n")
	}

	if sb.Len() > 0 {
		outPath := filepath.Join(typesDir, "go.d.ts")
		// Append to existing file if present
		existing, _ := os.ReadFile(outPath)
		_ = os.WriteFile(outPath, append(existing, []byte(sb.String())...), 0644)
	}
}
