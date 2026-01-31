package linker

import (
	"fmt"
	"os"
	"os/exec"
)

type Fetcher struct {
	TempDir string
}

func NewFetcher() (*Fetcher, error) {
	tempDir, err := os.MkdirTemp("", "typego-fetcher-*")
	if err != nil {
		return nil, err
	}

	// Initialize a dummy module to hold dependencies
	cmd := exec.Command("go", "mod", "init", "typego-dummy")
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to init dummy mod: %s", out)
	}

	return &Fetcher{TempDir: tempDir}, nil
}

func (f *Fetcher) Get(pkgPath string) error {
	cmd := exec.Command("go", "get", pkgPath)
	cmd.Dir = f.TempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go get failed: %s", out)
	}
	return nil
}

func (f *Fetcher) Cleanup() {
	os.RemoveAll(f.TempDir)
}
