package plugins

import (
	"fmt"
	"os"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/repyh/typego/internal/transformer/core"
	"github.com/repyh/typego/internal/transformer/visitors"
)

// DeferPlugin creates an esbuild plugin that applies the Defer transformation.
func DeferPlugin() api.Plugin {
	return api.Plugin{
		Name: "typego-defer",
		Setup: func(build api.PluginBuild) {

			// Initialize Visitors
			core.Visitors = nil
			core.RegisterVisitor(&visitors.DeferVisitor{})
			core.RegisterVisitor(&visitors.IotaVisitor{})

			// Broad filter to capture everything for debugging, then check extension manually
			build.OnLoad(api.OnLoadOptions{Filter: `.*`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				if !strings.HasSuffix(args.Path, ".ts") {
					return api.OnLoadResult{}, nil
				}

				source, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				// 2a. Convert TS -> JS (Preserve semantics, remove types)
				jsRes := api.Transform(string(source), api.TransformOptions{
					Loader: api.LoaderTS,
					Format: api.FormatCommonJS,
					Target: api.ES2015,
				})
				if len(jsRes.Errors) > 0 {
					return api.OnLoadResult{Errors: jsRes.Errors}, nil
				}
				jsCode := string(jsRes.Code)

				// 3. Apply Defer Transformer on the clean JS
				newCode, err := core.Transform(args.Path, jsCode)
				if err != nil {
					return api.OnLoadResult{
						Errors: []api.Message{{Text: fmt.Sprintf("transform error: %v", err)}},
					}, nil
				}

				// 4. Return to esbuild
				return api.OnLoadResult{
					Contents: &newCode,
					Loader:   api.LoaderJS,
				}, nil
			})
		},
	}
}
