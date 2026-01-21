package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/repyh/typego/pkg/cli/cmd"
	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add <module>[@version]",
	Short: "Add a Go module dependency to the project",
	Long: `Add a Go module to typego.modules.json without manual editing.

Examples:
  typego add github.com/gin-gonic/gin
  typego add github.com/gin-gonic/gin@v1.9.0
  typego add github.com/stretchr/testify@latest`,
	Args: cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to get cwd: %v", err))
			os.Exit(1)
		}

		arg := args[0]
		module := arg
		version := "latest"

		if idx := strings.LastIndex(arg, "@"); idx != -1 {
			module = arg[:idx]
			version = arg[idx+1:]
		}

		config, err := internal.ReadModulesConfig(cwd)
		if err != nil {
			internal.Error(fmt.Sprintf("Failed to read config: %v", err))
			internal.Info("Run 'typego init' first to create a project")
			os.Exit(1)
		}

		if existingVersion, ok := config.Dependencies[module]; ok {
			internal.Warn(fmt.Sprintf("%s already exists with version %s", module, existingVersion))
			internal.Info(fmt.Sprintf("Updating to %s", version))
		}

		if config.Dependencies == nil {
			config.Dependencies = make(map[string]string)
		}
		config.Dependencies[module] = version

		if err := internal.WriteModulesConfig(cwd, config); err != nil {
			internal.Error(fmt.Sprintf("Failed to write config: %v", err))
			os.Exit(1)
		}

		internal.Success(fmt.Sprintf("Added %s@%s", module, version))

		fmt.Println("📦 Regenerating types...")
		cmd.TypesCmd.SetArgs([]string{})
		if err := cmd.TypesCmd.Execute(); err != nil {
			internal.Warn("Types generation skipped (run 'typego install' to fetch and build)")
		}
	},
}
