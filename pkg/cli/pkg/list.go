package pkg

import (
	"fmt"
	"os"

	"github.com/repyh/typego/pkg/cli/internal"
	"github.com/spf13/cobra"
)

var listTree bool
var listJSON bool

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed Go module dependencies",
	Long: `Display all Go module dependencies defined in typego.modules.json.

Examples:
  typego list
  typego list --json`,
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

		if listJSON {
			fmt.Println("{")
			i := 0
			for mod, ver := range config.Dependencies {
				comma := ","
				if i == len(config.Dependencies)-1 {
					comma = ""
				}
				fmt.Printf("  %q: %q%s\n", mod, ver, comma)
				i++
			}
			fmt.Println("}")
			return
		}

		fmt.Printf("📦 Dependencies (%d):\n", len(config.Dependencies))
		for mod, ver := range config.Dependencies {
			fmt.Printf("  ├── %s@%s\n", mod, ver)
		}
	},
}

func init() {
	ListCmd.Flags().BoolVar(&listJSON, "json", false, "Output in JSON format")
	ListCmd.Flags().BoolVar(&listTree, "tree", false, "Show dependency tree (not implemented)")
}
