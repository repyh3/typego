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

var OutdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Show outdated dependencies",
	Long: `Check which dependencies have newer versions available.

Example:
  typego outdated`,
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
			internal.Info("No dependencies installed")
			return
		}

		internal.Step("🔍", "Checking for outdated dependencies...")

		outdatedCount := 0
		for module, currentVersion := range config.Dependencies {
			latestVersion, err := fetchLatestVersion(module)
			if err != nil {
				internal.Warn(fmt.Sprintf("Could not check %s: %v", module, err))
				continue
			}

			if currentVersion != latestVersion && currentVersion != "latest" {
				fmt.Printf("  📦 %s\n", module)
				fmt.Printf("     Current: %s -> Latest: %s\n", currentVersion, latestVersion)
				outdatedCount++
			}
		}

		if outdatedCount == 0 {
			internal.Success("All dependencies are up to date!")
		} else {
			fmt.Println()
			internal.Info(fmt.Sprintf("Run 'typego update' to update %d dependencies", outdatedCount))
		}
	},
}

func fetchLatestVersion(module string) (string, error) {
	tmpDir := filepath.Join(os.TempDir(), "typego-outdated-check")
	_ = os.MkdirAll(tmpDir, 0755)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	initCmd := exec.Command("go", "mod", "init", "temp")
	initCmd.Dir = tmpDir
	_ = initCmd.Run()

	cmd := exec.Command("go", "list", "-m", "-versions", module)
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return "unknown", nil
	}

	return parts[len(parts)-1], nil
}
