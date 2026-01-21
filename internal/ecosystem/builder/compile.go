package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CompileBinary builds the custom TypeGo binary from the generated main.go
func CompileBinary(buildDir, outputDir, outputName string) error {
	// target path
	exePath := filepath.Join(outputDir, outputName)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	cmd := exec.Command("go", "build", "-o", exePath, ".")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Optional: Support local dev mode replacement for testing
	// We might need to propagate this logic from the installer if we are developing TypeGo itself.

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}
	return nil
}
