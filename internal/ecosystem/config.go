package ecosystem

import (
	"fmt"
	"strings"
)

const ConfigFileName = "typego.modules.json"
const LockFileName = "typego.lock"

// Lockfile captures resolved versions for reproducible builds
type Lockfile struct {
	LockfileVersion int                        `json:"lockfileVersion"`
	Resolved        map[string]ResolvedPackage `json:"resolved"`
}

// ResolvedPackage contains the exact version resolved by go get
type ResolvedPackage struct {
	Version string `json:"version"`
}

func DefaultLockfile() Lockfile {
	return Lockfile{
		LockfileVersion: 1,
		Resolved:        make(map[string]ResolvedPackage),
	}
}

// ModuleConfig represents the schema for typego.modules.json
type ModuleConfig struct {
	Schema       string            `json:"$schema,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
	Replace      map[string]string `json:"replace,omitempty"`
	Compiler     CompilerConfig    `json:"compiler,omitempty"`
}

type CompilerConfig struct {
	GoVersion string   `json:"goVersion,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func DefaultConfig() ModuleConfig {
	return ModuleConfig{
		// placeholder for now
		Schema:       "https://typego.dev/schemas/modules.json",
		Dependencies: map[string]string{
			// Example dependency
			// "github.com/gin-gonic/gin": "v1.9.0",
		},
		Compiler: CompilerConfig{
			GoVersion: "1.24",
		},
	}
}

func (c *ModuleConfig) Validate() error {
	for dep := range c.Dependencies {
		if dep == "" {
			return fmt.Errorf("dependency path cannot be empty")
		}
		// Basic check: module paths usually contain a dot (domain) or are standard lib (which are handled elsewhere, but for go.mod valid paths usually have a dot)
		// For now, we just ensure it's not empty and doesn't contain spaces.
		if strings.Contains(dep, " ") {
			return fmt.Errorf("invalid dependency path %q: cannot contain spaces", dep)
		}
	}
	return nil
}
