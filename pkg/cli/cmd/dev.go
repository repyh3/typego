package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/repyh/typego/internal/ecosystem"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

var DevCmd = &cobra.Command{
	Use:   "dev [file]",
	Short: "Start development server with hot-reload",
	Long: `Start a development server that watches for file changes and automatically
restarts. Provides colored output and compilation timing.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		cwd, _ := os.Getwd()

		if ecosystem.IsHandoffRequired(cwd) {
			binaryPath, _ := ecosystem.GetJITBinaryPath(cwd)

			// Call the local binary with same args (dev [file])
			handoff := exec.Command(binaryPath, os.Args[1:]...)
			handoff.Stdout = os.Stdout
			handoff.Stderr = os.Stderr
			handoff.Stdin = os.Stdin
			handoff.Env = append(os.Environ(), ecosystem.HandoffEnvVar+"=true")

			if err := handoff.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					os.Exit(exitErr.ExitCode())
				}
				fmt.Printf("Handoff failed: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		absPath, err := filepath.Abs(filename)
		if err != nil {
			printError("Failed to resolve path: %v", err)
			return
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			printError("File not found: %s", filename)
			return
		}

		printBanner()
		printInfo("Watching %s", filepath.Base(filename))
		printInfo("Press Ctrl+C to stop")
		fmt.Println()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		runDevProcess(absPath)

		lastMod := getDevLastMod(absPath)
		ticker := time.NewTicker(150 * time.Millisecond) // Faster polling
		defer ticker.Stop()

		for {
			select {
			case <-sig:
				fmt.Println()
				printWarning("Shutting down...")
				killDevProcess()
				printSuccess("Development server stopped")
				return
			case <-ticker.C:
				currentMod := getDevLastMod(absPath)
				if currentMod.After(lastMod) {
					fmt.Println()
					printInfo("Change detected, restarting...")
					killDevProcess()
					runDevProcess(absPath)
					lastMod = currentMod
				}
			}
		}
	},
}

var devCurrentCmd *exec.Cmd

func killDevProcess() {
	if devCurrentCmd != nil && devCurrentCmd.Process != nil {
		_ = devCurrentCmd.Process.Kill()
		_ = devCurrentCmd.Wait()
		devCurrentCmd = nil
	}
}

func runDevProcess(file string) {
	start := time.Now()

	exe, _ := os.Executable()
	cmd := exec.Command(exe, "run", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		printError("Failed to start: %v", err)
		return
	}

	devCurrentCmd = cmd
	elapsed := time.Since(start)
	printSuccess("Started in %dms", elapsed.Milliseconds())

	go func() {
		err := cmd.Wait()
		if err != nil {
			// Only print if not killed intentionally
			if devCurrentCmd != nil {
				printError("Process exited: %v", err)
			}
		}
	}()
}

func getDevLastMod(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func printBanner() {
	fmt.Printf("%s╔══════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║      TypeGo Development Server       ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚══════════════════════════════════════╝%s\n", colorCyan, colorReset)
	fmt.Println()
}

func printInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[INFO]%s %s\n", colorBlue, colorReset, msg)
}

func printSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[OK]%s %s\n", colorGreen, colorReset, msg)
}

func printWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, msg)
}

func printError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, msg)
}

func init() {
}
