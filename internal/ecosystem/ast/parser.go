package ast

import (
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

// ParseFile parses a TypeScript/JavaScript file and returns its AST
func ParseFile(path string) (*api.BuildResult, string, error) {
	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	src := string(content)

	// Use esbuild's Transform API to just parse (or close enough)
	// Actually esbuild doesn't expose a raw "Parse" AST in the Go API easily without Build/Transform.
	// But we can use Build with 'Write: false' to get metadata if we needed imports,
	// but strictly speaking we want to visit the AST.
	//
	// However, the standard esbuild Go API is high-level. It doesn't expose the AST nodes directly
	// for walking in Go. It only exposes plugins.
	//
	// CRITICAL STRATEGY CHANGE:
	// Since esbuild Go API doesn't expose AST traversal, we will use a simpler REGEX approach
	// for the MVP, or rely on esbuild's "Metafile" output to find imports.
	//
	// Using Metafile is much more robust than regex and faster/safer.

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Bundle:      true,  // Bundle to trace imports
		Write:       false, // Don't write to disk
		Metafile:    true,  // Generate metadata about imports
		Platform:    api.PlatformNeutral,
		Target:      api.ESNext,
		LogLevel:    api.LogLevelSilent,

		// We mock all imports as external so esbuild doesn't fail trying to resolve them
		Plugins: []api.Plugin{
			{
				Name: "import-mock",
				Setup: func(build api.PluginBuild) {
					// Catch ALL imports and mark them as external so we just record them
					// BUT don't mark the entry point itself as external!
					build.OnResolve(api.OnResolveOptions{Filter: ".*"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						if args.Kind == api.ResolveEntryPoint {
							return api.OnResolveResult{}, nil
						}
						return api.OnResolveResult{External: true}, nil
					})
				},
			},
		},
	})

	if len(result.Errors) > 0 {
		return nil, "", fmt.Errorf("parse error: %s", result.Errors[0].Text)
	}

	return &result, src, nil
}
