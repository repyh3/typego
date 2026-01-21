package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/repyh/typego/internal/ecosystem"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the ecosystem cache",
	Long:  `Remove the .typego directory and all cached binaries/dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		path := filepath.Join(cwd, ecosystem.HiddenDirName)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("✨ Cache already clean")
			return
		}

		fmt.Printf("🧹 Cleaning %s...\n", path)
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("❌ Failed to clean: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✨ Cache wiped successfully")
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
