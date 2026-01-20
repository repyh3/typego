package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	typegoBinary string
	rootDir      string
)

func TestMain(m *testing.M) {
	// 1. Determine Root Dir
	_, filename, _, _ := runtime.Caller(0)
	rootDir = filepath.Join(filepath.Dir(filename), "..", "..")

	// 2. Build TypeGo CLI
	typegoBinary = filepath.Join(os.TempDir(), "typego_e2e.exe")
	if runtime.GOOS != "windows" {
		typegoBinary = filepath.Join(os.TempDir(), "typego_e2e")
	}

	fmt.Println("🚧 Building TypeGo CLI for E2E...")
	cmd := exec.Command("go", "build", "-o", typegoBinary, filepath.Join(rootDir, "cmd", "typego"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to build typego: %v\n", err)
		os.Exit(1)
	}

	// 3. Run Tests
	exitCode := m.Run()

	// 4. Cleanup
	os.Remove(typegoBinary)
	os.Exit(exitCode)
}

func TestCLI_Run_HelloWorld(t *testing.T) {
	examplePath := filepath.Join(rootDir, "examples", "01-hello-world.ts")

	cmd := exec.Command(typegoBinary, "run", examplePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	output := stdout.String()
	expected := "Hello from TypeGo TypeScript!"
	if !strings.Contains(output, expected) {
		t.Errorf("Output missing %q. Got:\n%s", expected, output)
	}

	expectedRuntime := "Runtime: Goja (Go-based JS Engine) - v1.0.0"
	if !strings.Contains(output, expectedRuntime) {
		t.Errorf("Output missing %q. Got:\n%s", expectedRuntime, output)
	}
}

func TestCLI_Build_HelloWorld(t *testing.T) {
	// Temporarily create output in temp dir
	tempDir := t.TempDir()
	examplePath := filepath.Join(rootDir, "examples", "01-hello-world.ts")
	outputPath := filepath.Join(tempDir, "hello.exe")
	if runtime.GOOS != "windows" {
		outputPath = filepath.Join(tempDir, "hello")
	}

	// 1. Build Compilation
	t.Logf("Running typego build...")
	cmd := exec.Command(typegoBinary, "build", "-o", outputPath, examplePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("Build failed: %v\nStderr: %s", err, stderr.String())
	}

	// 2. Verify Binary Exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Binary not found at %s", outputPath)
	}

	// 3. Execute Compiled Binary
	t.Logf("Executing compiled binary...")
	binCmd := exec.Command(outputPath)
	var stdout bytes.Buffer
	binCmd.Stdout = &stdout
	if err := binCmd.Run(); err != nil {
		t.Fatalf("Compiled binary execution failed: %v", err)
	}

	// 4. Verify Output
	output := stdout.String()
	expected := "Hello from TypeGo TypeScript!"
	if !strings.Contains(output, expected) {
		t.Errorf("Output missing %q. Got:\n%s", expected, output)
	}
}
