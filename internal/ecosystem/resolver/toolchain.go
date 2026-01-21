package resolver

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGoGet runs 'go get' in the specified directory for the given modules
func RunGoGet(dir string, modules []string) error {
	if len(modules) == 0 {
		return nil
	}

	args := append([]string{"get"}, modules...)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get failed: %w", err)
	}
	return nil
}

// RunGoModTidy runs 'go mod tidy' in the specified directory
func RunGoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}

// InitGoMod initializes a new go.mod file if it doesn't exist
func InitGoMod(dir, moduleName string) error {
	// check if go.mod exists
	if _, err := os.Stat(dir + "/go.mod"); err == nil {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}
	return nil
}
