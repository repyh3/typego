package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/repyh/typego/internal/ecosystem"
	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the ecosystem cache",
	Long:  `Remove the .typego directory and all cached binaries/dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		path := filepath.Join(cwd, ecosystem.HiddenDirName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			internal.Success("Cache already clean")
			return
		}

		internal.Step("🧹", fmt.Sprintf("Cleaning %s...", path))
		if err := os.RemoveAll(path); err != nil {
			internal.Error(fmt.Sprintf("Failed to clean: %v", err))
			os.Exit(1)
		}
		internal.Success("Cache wiped successfully")
	},
}
