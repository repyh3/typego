package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const tsConfigTemplate = `{
  "compilerOptions": {
    "target": "ESNext",
    "module": "ESNext",
    "moduleResolution": "Node",
    "esModuleInterop": true,
    "strict": true,
    "baseUrl": ".",
    "paths": {
      "go/*": [".typego/types/*"]
    },
    "skipLibCheck": true
  },
  "include": ["src/**/*", ".typego/types/**/*"]
}`

const indexTemplate = `import { Println } from "go/fmt";
import { Sleep } from "go/sync";

// You can now import NPM packages!
// import _ from "lodash";

async function main() {
    Println("🚀 TypeGo Project Initialized!");
    await Sleep(500);
    Println("Happy coding!");
}

main();
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new TypeGo project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing TypeGo project...")

		dirs := []string{"src", ".typego/types"}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error creating directory %s: %v\n", dir, err)
				return
			}
		}

		indexPath := filepath.Join("src", "index.ts")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			_ = os.WriteFile(indexPath, []byte(indexTemplate), 0644)
			fmt.Println("Created src/index.ts")
		}

		if _, err := os.Stat("tsconfig.json"); os.IsNotExist(err) {
			_ = os.WriteFile("tsconfig.json", []byte(tsConfigTemplate), 0644)
			fmt.Println("Created tsconfig.json")
		}

		// Initialize NPM
		if _, err := os.Stat("package.json"); os.IsNotExist(err) {
			fmt.Println("📦 Initializing NPM...")
			if err := execShellCmd("npm", "init", "-y"); err != nil {
				fmt.Printf("Warning: npm init failed: %v\n", err)
			} else {
				fmt.Println("📥 Installing @types/node...")
				if err := execShellCmd("npm", "install", "-D", "@types/node"); err != nil {
					fmt.Printf("Warning: failed to install types: %v\n", err)
				}
			}
		}

		dtsPath := filepath.Join(".typego", "types", "go.d.ts")

		found := false
		curr, _ := os.Getwd()
		for i := 0; i < 5; i++ {
			target := filepath.Join(curr, "go.d.ts")
			masterDts, err := os.ReadFile(target)
			if err == nil {
				_ = os.WriteFile(dtsPath, masterDts, 0644)
				fmt.Printf("Synced %s from %s\n", dtsPath, target)
				found = true
				break
			}
			curr = filepath.Dir(curr)
		}

		if !found {
			fmt.Println("Warning: master go.d.ts not found in search path. Please copy it manually to .typego/types/")
		}

		fmt.Println("\n✨ Project ready! Run with: typego run src/index.ts")
	},
}

func execShellCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(initCmd)
}
