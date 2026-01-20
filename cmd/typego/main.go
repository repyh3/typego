package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var memoryLimit uint64

var rootCmd = &cobra.Command{
	Use:   "typego",
	Short: "TypeGo is a TypeScript runtime for Go",
	Long: `A high-performance TypeScript runtime built on Go, enabling 
developers to harness Go's concurrency and memory efficiency while writing TypeScript.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().Uint64VarP(&memoryLimit, "memory-limit", "M", 128, "Memory limit for the JS engine in MB")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
