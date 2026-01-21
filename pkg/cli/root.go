package cli

import (
	"fmt"
	"os"

	"github.com/repyh/typego/pkg/cli/pkg"
	"github.com/spf13/cobra"
)

var MemoryLimit uint64

var RootCmd = &cobra.Command{
	Use:   "typego",
	Short: "TypeGo is a TypeScript runtime for Go",
	Long: `A high-performance TypeScript runtime built on Go, enabling 
developers to harness Go's concurrency and memory efficiency while writing TypeScript.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().Uint64VarP(&MemoryLimit, "memory-limit", "M", 128, "Memory limit for the JS engine in MB")

	// Package manager commands
	RootCmd.AddCommand(pkg.AddCmd)
	RootCmd.AddCommand(pkg.RemoveCmd)
	RootCmd.AddCommand(pkg.ListCmd)
	RootCmd.AddCommand(pkg.InstallCmd)
	RootCmd.AddCommand(pkg.CleanCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
