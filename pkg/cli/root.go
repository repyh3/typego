package cli

import (
	"fmt"
	"os"

	"github.com/repyh/typego/pkg/cli/cmd"
	"github.com/repyh/typego/pkg/cli/pkg"
	"github.com/spf13/cobra"
)

var MemoryLimit uint64

var RootCmd = &cobra.Command{
	Use:   "typego",
	Short: "TypeGo is a TypeScript runtime for Go",
	Long: `A high-performance TypeScript runtime built on Go, enabling 
developers to harness Go's concurrency and memory efficiency while writing TypeScript.`,
	Run: func(c *cobra.Command, args []string) {
		_ = c.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().Uint64VarP(&MemoryLimit, "memory-limit", "M", 128, "Memory limit for the JS engine in MB")

	// Core commands
	RootCmd.AddCommand(cmd.BuildCmd)
	RootCmd.AddCommand(cmd.DevCmd)
	RootCmd.AddCommand(cmd.RunCmd)
	RootCmd.AddCommand(cmd.InitCmd)
	RootCmd.AddCommand(cmd.TypesCmd)
	RootCmd.AddCommand(cmd.WatchCmd)

	// Package manager commands
	RootCmd.AddCommand(pkg.AddCmd)
	RootCmd.AddCommand(pkg.RemoveCmd)
	RootCmd.AddCommand(pkg.ListCmd)
	RootCmd.AddCommand(pkg.InstallCmd)
	RootCmd.AddCommand(pkg.CleanCmd)
	RootCmd.AddCommand(pkg.UpdateCmd)
	RootCmd.AddCommand(pkg.OutdatedCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
