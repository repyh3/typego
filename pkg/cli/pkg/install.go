package pkg

import (
	"fmt"
	"os"

	"github.com/repyh/typego/internal/ecosystem/installer"
	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install external Go modules for this project",
	Long: `Analyze your TypeScript source code, resolve 'go:*' imports, 
and build a project-specific TypeGo runtime with those modules linked in.

This creates a local binary in .typego/bin/typego-app.exe that is optimized for your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to get cwd: %v", err))
			os.Exit(1)
		}

		internal.Step("🚀", "Installing Ecosystem for "+cwd)
		if err := installer.RunInstall(cwd); err != nil {
			internal.Error(fmt.Sprintf("Installation failed: %v", err))
			os.Exit(1)
		}
	},
}
