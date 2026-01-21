package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var WatchCmd = &cobra.Command{
	Use:   "watch [file]",
	Short: "Run a file and restart on changes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		absPath, err := filepath.Abs(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("👀 Watching %s...\n", filepath.Base(filename))

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		// Initial Run
		runProcess(absPath)

		// Polling Loop (Simple & Robust)
		lastMod := getLastMod(absPath)
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-sig:
				fmt.Println("\nStopped watching.")
				killProcess()
				return
			case <-ticker.C:
				currentMod := getLastMod(absPath)
				if currentMod.After(lastMod) {
					fmt.Println("🔄 Change detected, restarting...")
					killProcess()
					runProcess(absPath)
					lastMod = currentMod
				}
			}
		}
	},
}

var currentCmd *exec.Cmd

func killProcess() {
	if currentCmd != nil && currentCmd.Process != nil {
		if runtimeOS := "windows"; runtimeOS == "windows" {
			_ = currentCmd.Process.Kill() // Force kill on Windows
		} else {
			_ = currentCmd.Process.Signal(syscall.SIGTERM)
		}
		_ = currentCmd.Wait()
		currentCmd = nil
	}
}

func runProcess(file string) {
	// Re-use "run" logic by invoking typego run subcommand
	// Note: In real dev, we might call internal function, but spawning subprocess ensures clean state
	// and handles panic/crashes without killing watcher.

	// Assuming 'typego' is in PATH or we use the current executable
	exe, _ := os.Executable()

	cmd := exec.Command(exe, "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start: %v\n", err)
		return
	}

	currentCmd = cmd

	// Wait in background to clean up zombies
	go func() {
		_ = cmd.Wait()
	}()
}

func getLastMod(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func init() {
	// Registered in root.go
}
