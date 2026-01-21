package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update [module]",
	Short: "Update dependencies to latest versions",
	Long: `Update all or specific Go module dependencies to their latest versions.

Examples:
  typego update                           # Update all dependencies
  typego update github.com/gin-gonic/gin  # Update specific module`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to get cwd: %v", err))
			os.Exit(1)
		}

		config, err := internal.ReadModulesConfig(cwd)
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to read config: %v", err))
			os.Exit(1)
		}

		if len(config.Dependencies) == 0 {
			internal.Info("No dependencies to update")
			return
		}

		// If specific module provided, update only that
		if len(args) > 0 {
			module := args[0]
			if _, ok := config.Dependencies[module]; !ok {
				internal.Error(fmt.Sprintf("%s is not in dependencies", module))
				os.Exit(1)
			}

			internal.Step("🔄", fmt.Sprintf("Updating %s...", module))
			latestVersion, err := getLatestVersion(module)
			if err != nil {
				internal.Warn(fmt.Sprintf("Could not fetch latest version: %v", err))
				latestVersion = "latest"
			}

			config.Dependencies[module] = latestVersion
			if err := internal.WriteModulesConfig(cwd, config); err != nil {
				internal.Error(fmt.Sprintf("Failed to write config: %v", err))
				os.Exit(1)
			}

			internal.Success(fmt.Sprintf("Updated %s to %s", module, latestVersion))
		} else {
			// Update all
			internal.Step("🔄", "Updating all dependencies...")
			updated := 0

			for module := range config.Dependencies {
				latestVersion, err := getLatestVersion(module)
				if err != nil {
					internal.Warn(fmt.Sprintf("Could not fetch %s: %v", module, err))
					continue
				}

				if config.Dependencies[module] != latestVersion {
					config.Dependencies[module] = latestVersion
					internal.Info(fmt.Sprintf("  %s -> %s", module, latestVersion))
					updated++
				}
			}

			if err := internal.WriteModulesConfig(cwd, config); err != nil {
				internal.Error(fmt.Sprintf("Failed to write config: %v", err))
				os.Exit(1)
			}

			if updated > 0 {
				internal.Success(fmt.Sprintf("Updated %d dependencies", updated))
			} else {
				internal.Info("All dependencies are up to date")
			}
		}

		internal.Info("Run 'typego install' to rebuild the JIT binary")
	},
}

// getLatestVersion queries go list to get the latest version of a module
func getLatestVersion(module string) (string, error) {
	// Create temp dir for go list
	tmpDir := filepath.Join(os.TempDir(), "typego-update-check")
	_ = os.MkdirAll(tmpDir, 0755)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Initialize a dummy go.mod
	initCmd := exec.Command("go", "mod", "init", "temp")
	initCmd.Dir = tmpDir
	_ = initCmd.Run()

	// Query latest version
	cmd := exec.Command("go", "list", "-m", "-versions", module)
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse output: "module/path v1.0.0 v1.1.0 v1.2.0"
	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return "latest", nil
	}

	// Return the last version (most recent)
	return parts[len(parts)-1], nil
}

func init() {
	// Add to subcommand of typego if needed
}
