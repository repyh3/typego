package internal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/repyh/typego/internal/ecosystem"
)

// ReadModulesConfig reads and parses the typego.modules.json file
func ReadModulesConfig(cwd string) (*ecosystem.ModuleConfig, error) {
	configPath := filepath.Join(cwd, ecosystem.ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ecosystem.ModuleConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// WriteModulesConfig writes the config back to typego.modules.json
func WriteModulesConfig(cwd string, cfg *ecosystem.ModuleConfig) error {
	configPath := filepath.Join(cwd, ecosystem.ConfigFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// FindProjectRoot walks up from startDir to find a typego.modules.json
func FindProjectRoot(startDir string) (string, bool) {
	curr := startDir
	for {
		configPath := filepath.Join(curr, ecosystem.ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			return curr, true
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", false
}
