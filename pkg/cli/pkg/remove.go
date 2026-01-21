package pkg

import (
	"fmt"
	"os"

	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:     "remove <module>",
	Aliases: []string{"rm", "uninstall"},
	Short:   "Remove a Go module dependency from the project",
	Long: `Remove a Go module from typego.modules.json.

Examples:
  typego remove github.com/gin-gonic/gin
  typego rm github.com/gin-gonic/gin`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to get cwd: %v", err))
			os.Exit(1)
		}

		module := args[0]

		config, err := internal.ReadModulesConfig(cwd)
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to read config: %v", err))
			os.Exit(1)
		}

		if _, ok := config.Dependencies[module]; !ok {
			internal.Warn(fmt.Sprintf("%s is not in dependencies", module))
			os.Exit(0)
		}

		delete(config.Dependencies, module)

		if err := internal.WriteModulesConfig(cwd, config); err != nil {
			internal.Error(fmt.Sprintf("Failed to write config: %v", err))
			os.Exit(1)
		}

		internal.Success(fmt.Sprintf("Removed %s", module))
		internal.Info("Run 'typego install' to rebuild the JIT binary")
	},
}
