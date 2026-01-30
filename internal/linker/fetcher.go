package linker

import (
	"fmt"
	"os"
	"os/exec"
)

// Fetcher handles downloading remote Go modules
type Fetcher struct {
	TempDir string
}

// NewFetcher creates a new Fetcher in a temporary directory
func NewFetcher() (*Fetcher, error) {
	tempDir, err := os.MkdirTemp("", "typego-fetcher-*")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "mod", "init", "typego-dummy")
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to init dummy mod: %s", out)
	}

	return &Fetcher{TempDir: tempDir}, nil
}

// Get downloading a package version using 'go get'
func (f *Fetcher) Get(pkgPath string) error {
	cmd := exec.Command("go", "get", pkgPath)
	cmd.Dir = f.TempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go get failed: %s", out)
	}
	return nil
}

// Cleanup removes the temporary directory
func (f *Fetcher) Cleanup() {
	os.RemoveAll(f.TempDir)
}
