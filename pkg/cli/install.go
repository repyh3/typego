package cli

import (
	"fmt"
	"os"

	"github.com/repyh/typego/internal/ecosystem/installer"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install external Go modules for this project",
	Long: `Analyze your TypeScript source code, resolve 'go:*' imports, 
and build a project-specific TypeGo runtime with those modules linked in.

This creates a local binary in .typego/bin/typego-app.exe that is optimized for your project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting cwd: %v\n", err)
			return
		}

		fmt.Println("🚀 Installing Ecosystem for " + cwd)
		if err := installer.RunInstall(cwd); err != nil {
			fmt.Printf("❌ Installation failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
