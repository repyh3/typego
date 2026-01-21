package ecosystem

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

const (
	HiddenDirName = ".typego"
	BinaryName    = "typego-app.exe"
	HandoffEnvVar = "TYPEGO_HANDOFF"
)

// GetJITBinaryPath returns the path to the JIT binary if it exists
func GetJITBinaryPath(cwd string) (string, bool) {
	path := filepath.Join(cwd, HiddenDirName, "bin", BinaryName)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

// IsHandoffRequired returns true if we should delegate to the JIT binary
func IsHandoffRequired(cwd string) bool {
	// Prevent infinite loops where the JIT binary calls itself
	if os.Getenv(HandoffEnvVar) == "true" {
		return false
	}

	_, exists := GetJITBinaryPath(cwd)
	if !exists {
		return false
	}

	// Stale check
	return VerifyChecksum(cwd)
}

// EnsureGitIgnore adds HiddenDirName to .gitignore if it exists and doesn't already contain it
func EnsureGitIgnore(cwd string) error {
	gitIgnorePath := filepath.Join(cwd, ".gitignore")

	// Read existing content if it exists
	content, err := os.ReadFile(gitIgnorePath)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == HiddenDirName || strings.TrimSpace(line) == HiddenDirName+"/" {
				return nil // Already ignored
			}
		}
	}

	// Create or Append to file
	f, err := os.OpenFile(gitIgnorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(content) > 0 && content[len(content)-1] != '\n' {
		_, _ = f.WriteString("\n")
	}
	_, err = f.WriteString(HiddenDirName + "/\n")
	return err
}

// GetConfigHash returns the SHA256 hash of the typego.modules.json file
func GetConfigHash(cwd string) (string, error) {
	configPath := filepath.Join(cwd, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // No config, no hash
		}
		return "", err
	}
	// Use MD5 or SHA256, sha256 is better
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// WriteChecksum saves the hash to .typego/checksum
func WriteChecksum(cwd string, hash string) error {
	path := filepath.Join(cwd, HiddenDirName, "checksum")
	return os.WriteFile(path, []byte(hash), 0644)
}

// VerifyChecksum returns true if the current config hash matches the saved checksum
func VerifyChecksum(cwd string) bool {
	savedPath := filepath.Join(cwd, HiddenDirName, "checksum")
	saved, err := os.ReadFile(savedPath)
	if err != nil {
		return false
	}

	current, err := GetConfigHash(cwd)
	if err != nil {
		return false
	}

	return string(saved) == current
}

// FindRepoRoot attempts to find the root of the typego repository by walking up from startDir
func FindRepoRoot(startDir string) (string, bool) {
	curr := startDir
	for {
		goModPath := filepath.Join(curr, "go.mod")
		if data, err := os.ReadFile(goModPath); err == nil {
			if strings.Contains(string(data), "module github.com/repyh/typego") {
				return curr, true
			}
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", false
}
