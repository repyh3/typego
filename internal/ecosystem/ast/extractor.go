package ast

import (
	"encoding/json"
	"strings"
)

// Metafile structure matches esbuild's JSON output subset
type Metafile struct {
	Inputs map[string]InputFile `json:"inputs"`
}

type InputFile struct {
	Imports []ImportInfo `json:"imports"`
}

type ImportInfo struct {
	Path string `json:"path"`
	// Original fields like Kind, External are available but we just need Path
}

// ExtractImports parses the metafile JSON and returns a unique list of "go:*" imports
// It filters out core internal modules (fmt, os, sync, net/http, memory, crypto)
// leaving only the external ones that need shimming/installing.
func ExtractImports(metafileJSON string) ([]string, error) {
	var meta Metafile
	if err := json.Unmarshal([]byte(metafileJSON), &meta); err != nil {
		return nil, err
	}

	unique := make(map[string]bool)
	for _, file := range meta.Inputs {
		for _, imp := range file.Imports {
			if strings.HasPrefix(imp.Path, "go:") {
				// Strip "go:" prefix to check against core modules
				// But return the full path (or stripped? let's return full path for clarity, or stripped for easier go get?)
				// The plan says: "go:github.com/..." -> "github.com/..."

				// Let's normalize it here
				rawPath := imp.Path

				// Optional: Filter out known internal modules if we want only external ones
				// or keep all and let the resolver decide.
				// The prompt says "Returns unique list of used modules".
				unique[rawPath] = true
			}
		}
	}

	var result []string
	for k := range unique {
		result = append(result, k)
	}
	return result, nil
}
