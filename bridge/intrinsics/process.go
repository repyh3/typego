package intrinsics

import (
	"os"
	"runtime"
	"strings"

	"github.com/grafana/sobek"
)

// EnableProcess injects the Node.js `process` global
func (r *Registry) EnableProcess() {
	proc := r.vm.NewObject()

	env := r.vm.NewObject()
	whitelist := map[string]bool{
		"PATH":     true,
		"LANG":     true,
		"PWD":      true,
		"HOSTNAME": true,
		"USER":     true,
	}

	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			upperKey := strings.ToUpper(key)
			// Allow whitelisted vars or anything prefixed with TYPEGO_
			if whitelist[upperKey] || strings.HasPrefix(upperKey, "TYPEGO_") {
				_ = env.Set(key, parts[1])
			}
		}
	}

	_ = env.Set("FORCE_COLOR", "1")
	_ = proc.Set("env", env)

	_ = proc.Set("platform", runtime.GOOS)

	_ = proc.Set("cwd", func(call sobek.FunctionCall) sobek.Value {
		wd, _ := os.Getwd()
		return r.vm.ToValue(wd)
	})

	_ = proc.Set("argv", os.Args)

	_ = proc.Set("version", runtime.Version())

	_ = r.vm.Set("process", proc)
}
