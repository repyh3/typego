package dx_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestTypeGeneration(t *testing.T) {
	// Create a temporary workspace
	tmpDir := t.TempDir()

	// Create a reproduction file that uses standard libraries
	reproFile := filepath.Join(tmpDir, "main.ts")
	reproContent := `
import * as fmt from "go:fmt";
import * as os from "go:os";
import * as url from "go:net/url";
import * as http from "go:net/http";
import { makeShared } from "typego:memory";

export function main() {
	fmt.Println("Hello");
	const u = url.Parse("https://example.com");
	fmt.Println(u.Host);
}
`
	if err := os.WriteFile(reproFile, []byte(reproContent), 0644); err != nil {
		t.Fatalf("Failed to write repro file: %v", err)
	}

	// Initialize TypeGo project structure
	if err := os.MkdirAll(filepath.Join(tmpDir, ".typego", "types"), 0755); err != nil {
		t.Fatalf("Failed to create .typego dir: %v", err)
	}

	// Run 'typego types' relative to the package directory
	cwd, _ := os.Getwd()
	rootDir := filepath.Join(cwd, "../../")

	// Build the binary first to avoid go run module issues in tmpDir
	binPath := filepath.Join(tmpDir, "typego")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/typego")
	buildCmd.Dir = rootDir
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build typego: %s", out)
	}

	cmd := exec.Command(binPath, "types", "main.ts")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("typego types failed: %s", out)
	}

	// Check go.d.ts content
	dtsPath := filepath.Join(tmpDir, ".typego", "types", "go.d.ts")
	contentBytes, err := os.ReadFile(dtsPath)
	if err != nil {
		t.Fatalf("Failed to read go.d.ts: %v", err)
	}
	content := string(contentBytes)

	// Assertions
	checks := []string{
		`declare module "go:fmt"`,
		`declare module "go:os"`,
		`declare module "go:net/url"`,
		`declare module "go:net/http"`,
		`declare module "typego:memory"`,
		`function isGoError(e: unknown): e is Error;`, // Global helper
		`Parse(rawURL: string): URL;`,                 // Ensure net/url Parse exists
	}

	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("go.d.ts missing expected content: %s", check)
		}
	}

	// Ensure NO wildcard
	if strings.Contains(content, `declare module "go:*"`) {
		t.Errorf("go.d.ts should NOT contain wildcard go:*")
	}
}
