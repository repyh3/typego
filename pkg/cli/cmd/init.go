package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initWithNpm bool

const tsConfigTemplate = `{
  "compilerOptions": {
    "target": "ESNext",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "esModuleInterop": true,
    "strict": true,
    "baseUrl": ".",
    "paths": {
      "go:*": [".typego/types/*"]
    },
    "skipLibCheck": true
  },
  "include": ["src/**/*", ".typego/types/**/*"]
}`

const indexTemplate = `import { Println } from "go:fmt";
import { Sleep } from "go:sync";

async function main() {
    Println("🚀 TypeGo Project Initialized!");
    await Sleep(500);
    Println("Happy coding with Go packages!");
}

main();
`

const modulesTemplate = `{
  "dependencies": {},
  "compiler": {
    "goVersion": "1.24"
  }
}`

var InitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new TypeGo project",
	Long: `Initialize a new TypeGo project with Go-first defaults.

By default, npm is NOT initialized. Use --npm flag if you need npm packages.
TypeScript type definitions are automatically generated.

Examples:
  typego init              # Initialize in current directory
  typego init my-app       # Create my-app directory and initialize
  typego init --npm        # Include npm/package.json setup`,
	Run: func(cmd *cobra.Command, args []string) {
		projectDir := "."
		if len(args) > 0 {
			projectDir = args[0]
			if err := os.MkdirAll(projectDir, 0755); err != nil {
				fmt.Printf("❌ Error creating project directory: %v\n", err)
				return
			}
		}

		fmt.Println("🚀 Creating TypeGo project...")
		fmt.Println()

		dirs := []string{
			filepath.Join(projectDir, "src"),
			filepath.Join(projectDir, ".typego", "types"),
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("❌ Error creating directory %s: %v\n", dir, err)
				return
			}
		}

		files := []struct {
			path    string
			content string
			name    string
		}{
			{filepath.Join(projectDir, "src", "index.ts"), indexTemplate, "src/index.ts"},
			{filepath.Join(projectDir, "typego.modules.json"), modulesTemplate, "typego.modules.json"},
			{filepath.Join(projectDir, "tsconfig.json"), tsConfigTemplate, "tsconfig.json"},
		}

		fmt.Println("📁 Created project structure:")
		for _, f := range files {
			if _, err := os.Stat(f.path); os.IsNotExist(err) {
				if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
					fmt.Printf("❌ Error creating %s: %v\n", f.name, err)
					return
				}
				fmt.Printf("   ├── %s\n", f.name)
			} else {
				fmt.Printf("   ├── %s (exists, skipped)\n", f.name)
			}
		}

		fmt.Println("   └── .typego/types/go.d.ts")
		typesCmd := exec.Command("typego", "types")
		typesCmd.Dir = projectDir
		if err := typesCmd.Run(); err != nil {
			fmt.Println("   ⚠️  Could not auto-generate types. Run 'typego types' manually.")
		}

		if initWithNpm {
			fmt.Println()
			fmt.Println("📦 Initializing npm (--npm flag)...")
			if _, err := os.Stat(filepath.Join(projectDir, "package.json")); os.IsNotExist(err) {
				npmInit := exec.Command("npm", "init", "-y")
				npmInit.Dir = projectDir
				if err := npmInit.Run(); err != nil {
					fmt.Printf("   ⚠️  npm init failed: %v\n", err)
				} else {
					npmInstall := exec.Command("npm", "install", "-D", "@types/node")
					npmInstall.Dir = projectDir
					if err := npmInstall.Run(); err != nil {
						fmt.Printf("   ⚠️  failed to install @types/node: %v\n", err)
					} else {
						fmt.Println("   ✅ Installed @types/node")
					}
				}
			}
		}

		fmt.Println()
		fmt.Println("📋 Next steps:")
		if projectDir != "." {
			fmt.Printf("   1. cd %s\n", projectDir)
			fmt.Println("   2. typego run src/index.ts")
		} else {
			fmt.Println("   1. typego run src/index.ts")
		}
		fmt.Println()
		fmt.Println("💡 Tip: Add Go dependencies with 'typego add github.com/gin-gonic/gin'")
	},
}

func init() {
	InitCmd.Flags().BoolVar(&initWithNpm, "npm", false, "Initialize npm and install @types/node")
}
