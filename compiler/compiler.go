package compiler

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/repyh/typego/compiler/plugins"
)

type Result struct {
	JS        string
	SourceMap string
	Imports   []string
}

// GlobalVirtualModules allow pre-registering modules for JIT binaries
var GlobalVirtualModules = make(map[string]string)

func Compile(entryPoint string, virtualModules map[string]string) (*Result, error) {
	if virtualModules == nil {
		virtualModules = make(map[string]string)
	}
	for k, v := range GlobalVirtualModules {
		if _, exists := virtualModules[k]; !exists {
			virtualModules[k] = v
		}
	}

	if len(virtualModules) == 0 {
		if res, err := CheckCache(entryPoint); err == nil && res != nil {
			return res, nil
		}
	}

	var collectedImports []string

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryPoint},
		Bundle:      true,
		Write:       false,
		LogLevel:    api.LogLevelSilent,
		Target:      api.ESNext,
		Format:      api.FormatIIFE,
		Sourcemap:   api.SourceMapInline,
		Plugins: []api.Plugin{
			plugins.DeferPlugin(),
			{
				Name: "typego-virtual",
				Setup: func(build api.PluginBuild) {
					build.OnResolve(api.OnResolveOptions{Filter: `^go:.*`}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						collectedImports = append(collectedImports, args.Path)

						switch args.Path {
						case "go:fmt", "go:os", "go:sync", "go:net/http", "go:memory", "go:crypto":
							return api.OnResolveResult{Path: args.Path, Namespace: "typego-internal"}, nil
						}

						return api.OnResolveResult{Path: args.Path, Namespace: "typego-hyperlink"}, nil
					})
					build.OnResolve(api.OnResolveOptions{Filter: `^typego:.*`}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						collectedImports = append(collectedImports, args.Path)
						return api.OnResolveResult{Path: args.Path, Namespace: "typego-internal"}, nil
					})
					build.OnResolve(api.OnResolveOptions{Filter: `^go/.*`}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						collectedImports = append(collectedImports, args.Path)
						return api.OnResolveResult{Path: args.Path, Namespace: "typego-internal"}, nil
					})
					build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "typego-hyperlink"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						if content, ok := virtualModules[args.Path]; ok {
							return api.OnLoadResult{Contents: &content, Loader: api.LoaderTS}, nil
						}

						// Default/Fallback (used during scan)
						content := "const p = new Proxy({}, { get: () => () => {} }); export const Println = p; export const Printf = p; export default p;"
						return api.OnLoadResult{Contents: &content, Loader: api.LoaderTS}, nil
					})
					build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "typego-internal"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						var content string
						switch args.Path {
						case "go:memory":
							content = "const m = (globalThis as any).__typego_memory__; export const Ptr = m.ptr; export const makeShared = m.makeShared;"
						case "go:fmt":
							content = "const f = (globalThis as any).__go_fmt__; export const Println = f.Println; export const Printf = f.Printf;"
						case "go:os":
							content = "const o = (globalThis as any).__go_os__; export const WriteFile = o.WriteFile; export const ReadFile = o.ReadFile;"
						case "go:net/http":
							content = "const h = (globalThis as any).__go_http__; export const Get = h.Get; export const Fetch = h.Fetch; export const Post = h.Post; export const ListenAndServe = h.ListenAndServe;"
						case "go:sync":
							content = "const s = (globalThis as any).__go_sync__; export const Spawn = s.Spawn; export const Sleep = s.Sleep; export const Chan = (globalThis as any).Chan;"
						case "go:crypto":
							content = "const c = (globalThis as any).__go_crypto__; export const Sha256 = c.Sha256; export const Sha512 = c.Sha512; export const HmacSha256 = c.HmacSha256; export const HmacSha256Verify = c.HmacSha256Verify; export const RandomBytes = c.RandomBytes; export const Uuid = c.Uuid;"

						case "typego:memory":
							content = "const m = (globalThis as any).__typego_memory__; export const makeShared = m.makeShared; export const stats = m.stats; export const ptr = m.ptr;"
						case "typego:worker":
							content = "const w = (globalThis as any).__typego_worker__; export const Worker = w.Worker;"
						default:
							return api.OnLoadResult{Errors: []api.Message{{Text: "Unknown virtual module: " + args.Path}}}, nil
						}
						return api.OnLoadResult{Contents: &content, Loader: api.LoaderTS}, nil
					})
				},
			},
		},
	})

	res := &Result{
		Imports: collectedImports,
	}

	if len(result.Errors) > 0 {
		return res, fmt.Errorf("compilation failed: %v", result.Errors[0].Text)
	}

	for _, file := range result.OutputFiles {
		if file.Path == "<stdout>.js" || len(result.OutputFiles) == 1 {
			res.JS = string(file.Contents)
		} else if file.Path == "<stdout>.js.map" {
			res.SourceMap = string(file.Contents)
		}
	}

	if res.JS == "" && len(result.OutputFiles) > 0 {
		res.JS = string(result.OutputFiles[0].Contents)
	}

	if len(virtualModules) == 0 {
		_ = SaveCache(entryPoint, res)
	}

	return res, nil
}
